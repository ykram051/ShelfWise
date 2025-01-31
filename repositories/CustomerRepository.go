package repositories

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"FinalProject/models"
)

type CustomerStore interface {
	CreateCustomer(c models.Customer) (models.Customer, error)
	GetCustomer(id int) (models.Customer, error)
	UpdateCustomer(id int, c models.Customer) (models.Customer, error)
	DeleteCustomer(id int) error
	ListCustomers() ([]models.Customer, error)
	Save() error
}

type InMemoryCustomerStore struct {
	customers map[int]models.Customer
	nextID    int
	mu        sync.Mutex
	backend   string
}

func NewInMemoryCustomerStore(fileName string) *InMemoryCustomerStore {
	store := &InMemoryCustomerStore{
		customers: make(map[int]models.Customer),
		backend:   fileName,
	}
	LoadFromFile(fileName, &store.customers, &store.mu)

	for id := range store.customers {
		if id > store.nextID {
			store.nextID = id
		}
	}
	return store
}

func (s *InMemoryCustomerStore) CreateCustomer(c models.Customer) (models.Customer, error) {
	s.mu.Lock()
	s.nextID++
	c.ID = s.nextID
	c.CreatedAt = time.Now().UTC()
	s.customers[c.ID] = c
	s.mu.Unlock()
	err := s.Save()
	if err != nil {
		return models.Customer{}, fmt.Errorf("error saving customers: %v", err)
	}
	return c, nil
}

func (s *InMemoryCustomerStore) GetCustomer(id int) (models.Customer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.customers[id]
	if !ok {
		return models.Customer{}, errors.New("customer not found")
	}
	return c, nil
}

func (s *InMemoryCustomerStore) UpdateCustomer(id int, c models.Customer) (models.Customer, error) {
	s.mu.Lock()
	_, ok := s.customers[id]
	if !ok {
		return models.Customer{}, errors.New("customer not found")
	}
	c.ID = id
	s.customers[id] = c
	s.mu.Unlock()
	err := s.Save()
	if err != nil {
		return models.Customer{}, fmt.Errorf("error saving customers: %v", err)
	}
	return c, nil
}

func (s *InMemoryCustomerStore) DeleteCustomer(id int) error {
	s.mu.Lock()
	if _, ok := s.customers[id]; !ok {
		return errors.New("customer not found")
	}
	delete(s.customers, id)
	s.mu.Unlock()
	err := s.Save()
	if err != nil {
		return fmt.Errorf("error saving customers: %v", err)
	}
	return nil
}

func (s *InMemoryCustomerStore) ListCustomers() ([]models.Customer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []models.Customer
	for _, c := range s.customers {
		result = append(result, c)
	}
	return result, nil
}

func (s *InMemoryCustomerStore) Save() error {
	if err := SaveToFile(s.backend, s.customers, &s.mu); err != nil {
		return fmt.Errorf("error saving customers: %v", err)
	}
	return nil
}
