package repositories

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"FinalProject/models"
)

type OrderStore interface {
	CreateOrder(o models.Order) (models.Order, error)
	GetOrder(id int) (models.Order, error)
	UpdateOrder(id int, o models.Order) (models.Order, error)
	DeleteOrder(id int) error
	ListOrders() ([]models.Order, error)
	GetOrdersByDateRange(from, to time.Time) ([]models.Order, error)
	Save() error
}

type InMemoryOrderStore struct {
	orders  map[int]models.Order
	nextID  int
	mu      sync.Mutex
	backend string
}

func NewInMemoryOrderStore(fileName string) *InMemoryOrderStore {
	store := &InMemoryOrderStore{
		orders:  make(map[int]models.Order),
		backend: fileName,
	}
	LoadFromFile(fileName, &store.orders, &store.mu)

	for id := range store.orders {
		if id > store.nextID {
			store.nextID = id
		}
	}
	return store
}

func (s *InMemoryOrderStore) CreateOrder(o models.Order) (models.Order, error) {
	s.mu.Lock()
	s.nextID++
	o.ID = s.nextID
	o.CreatedAt = time.Now().UTC()
	var total float64
	for _, item := range o.Items {
		total += item.Book.Price * float64(item.Quantity)
	}
	o.TotalPrice = total
	o.Status = "Created"

	s.orders[o.ID] = o
	s.mu.Unlock()
	err := s.Save()
	if err != nil {
		return models.Order{}, fmt.Errorf("error saving orders: %v", err)
	}
	return o, nil
}

func (s *InMemoryOrderStore) GetOrder(id int) (models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[id]
	if !ok {
		return models.Order{}, errors.New("order not found")
	}
	return order, nil
}

func (s *InMemoryOrderStore) UpdateOrder(id int, o models.Order) (models.Order, error) {
	s.mu.Lock()
	_, ok := s.orders[id]
	if !ok {
		return models.Order{}, errors.New("order not found")
	}
	o.ID = id
	var total float64
	for _, item := range o.Items {
		total += item.Book.Price * float64(item.Quantity)
	}
	o.TotalPrice = total
	s.orders[id] = o
	s.mu.Unlock()
	err := s.Save()
	if err != nil {
		return models.Order{}, fmt.Errorf("error saving orders: %v", err)
	}
	return o, nil
}

func (s *InMemoryOrderStore) DeleteOrder(id int) error {
	s.mu.Lock()
	if _, ok := s.orders[id]; !ok {
		return errors.New("order not found")
	}
	delete(s.orders, id)
	s.mu.Unlock()

	err := s.Save()
	if err != nil {
		return fmt.Errorf("error saving orders: %v", err)
	}
	return nil
}

func (s *InMemoryOrderStore) ListOrders() ([]models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []models.Order
	for _, o := range s.orders {
		result = append(result, o)
	}
	return result, nil
}

func (s *InMemoryOrderStore) GetOrdersByDateRange(from, to time.Time) ([]models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var results []models.Order
	for _, o := range s.orders {
		if (o.CreatedAt.Equal(from) || o.CreatedAt.After(from)) &&
			(o.CreatedAt.Equal(to) || o.CreatedAt.Before(to)) {
			results = append(results, o)
		}
	}
	return results, nil
}

func (s *InMemoryOrderStore) Save() error {
	if err := SaveToFile(s.backend, s.orders, &s.mu); err != nil {
		return fmt.Errorf("error saving orders: %v", err)
	}
	return nil
}
