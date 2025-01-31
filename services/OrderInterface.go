package services

import (
	"FinalProject/models"
	"FinalProject/repositories"
	"context"
	"fmt"
	"time"
)

type OrderService struct {
	store     repositories.OrderStore
	bookstore repositories.BookStore
	customerstore repositories.CustomerStore
}

func NewOrderService(store repositories.OrderStore, bookstore repositories.BookStore,customerstore repositories.CustomerStore) *OrderService {
	return &OrderService{store: store,
		bookstore: bookstore,
		customerstore: customerstore}
}

func (s *OrderService) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	select {
	case <-ctx.Done():
		return models.Order{}, ctx.Err()
	default:
		customer, err := s.customerstore.GetCustomer(order.Customer.ID)
		if err != nil {
			return models.Order{}, fmt.Errorf("customer with ID %d not found: %v", order.Customer.ID, err)
		}
		order.Customer = customer

		for i, item := range order.Items {
			book, err := s.bookstore.GetBook(item.Book.ID)
			if err != nil {
				return models.Order{}, err
			}

			if book.Stock < item.Quantity {
				return models.Order{}, fmt.Errorf("insufficient stock for book ID %d", item.Book.ID)
			}

			book.Stock -= item.Quantity
			if _, err := s.bookstore.UpdateBook(item.Book.ID, book); err != nil {
				return models.Order{}, err
			}

			order.Items[i].Book = book
		}

		createdOrder, err := s.store.CreateOrder(order)
		if err != nil {
			return models.Order{}, err
		}

		return createdOrder, nil
	}
}

func (s *OrderService) GetOrder(ctx context.Context, id int) (models.Order, error) {
	select {
	case <-ctx.Done():
		return models.Order{}, ctx.Err()
	default:
	}
	return s.store.GetOrder(id)
}

func (s *OrderService) UpdateOrder(ctx context.Context, id int, updatedOrder models.Order) (models.Order, error) {
	select {
	case <-ctx.Done():
		return models.Order{}, ctx.Err()
	default:
		existingOrder, err := s.store.GetOrder(id)
		if err != nil {
			return models.Order{}, err
		}

		for _, item := range existingOrder.Items {
			book, err := s.bookstore.GetBook(item.Book.ID)
			if err != nil {
				return models.Order{}, err
			}
			book.Stock += item.Quantity
			if _, err := s.bookstore.UpdateBook(item.Book.ID, book); err != nil {
				return models.Order{}, err
			}
		}

		for i, item := range updatedOrder.Items {
			book, err := s.bookstore.GetBook(item.Book.ID)
			if err != nil {
				return models.Order{}, err
			}

			if book.Stock < item.Quantity {
				return models.Order{}, fmt.Errorf("insufficient stock for book ID %d", item.Book.ID)
			}

			book.Stock -= item.Quantity
			if _, err := s.bookstore.UpdateBook(item.Book.ID, book); err != nil {
				return models.Order{}, err
			}

			updatedOrder.Items[i].Book = book
		}

		updatedOrder, err = s.store.UpdateOrder(id, updatedOrder)
		if err != nil {
			return models.Order{}, err
		}

		return updatedOrder, nil
	}
}

func (s *OrderService) DeleteOrder(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return s.store.DeleteOrder(id)
}

func (s *OrderService) ListOrders(ctx context.Context) ([]models.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.store.ListOrders()
}

func (s *OrderService) GetOrdersInRange(ctx context.Context, from, to time.Time) ([]models.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.store.GetOrdersByDateRange(from, to)
}

func (s *OrderService) SaveChanges() error {
	return s.store.Save()
}
