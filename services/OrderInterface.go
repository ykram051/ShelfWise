package services

import (
	"FinalProject/models"
	"FinalProject/repositories"
	"context"
	"fmt"
	"log"
	"time"
)

type OrderService struct {
	store         repositories.OrderStore
	bookstore     repositories.BookStore
	customerstore repositories.CustomerStore
}

func NewOrderService(store repositories.OrderStore, bookstore repositories.BookStore, customerstore repositories.CustomerStore) *OrderService {
	return &OrderService{store: store, bookstore: bookstore, customerstore: customerstore}
}

// CreateOrder processes an order with stock updates
func (s *OrderService) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	select {
	case <-ctx.Done():
		return models.Order{}, ctx.Err()
	default:
	}

	customer, err := s.customerstore.GetCustomer(order.CustomerID)
	if err != nil {
		return models.Order{}, fmt.Errorf("customer with ID %d not found: %w", order.CustomerID, err)
	}
	order.Customer = &customer  // ✅ Store Customer

	var total float64
	for i, item := range order.Items {
		book, err := s.bookstore.GetBook(item.BookID)
		if err != nil {
			return models.Order{}, fmt.Errorf("book with ID %d not found: %w", item.BookID, err)
		}

		if book.Stock < item.Quantity {
			return models.Order{}, fmt.Errorf("insufficient stock for book ID %d", item.BookID)
		}

		// ✅ Deduct stock and update book
		book.Stock -= item.Quantity
		if _, err := s.bookstore.UpdateBook(book.ID, book); err != nil {
			return models.Order{}, err
		}

		// ✅ Calculate total price
		total += float64(item.Quantity) * book.Price

		// ✅ Store the book inside the order items
		order.Items[i].Book = &book
	}

	order.TotalPrice = total  // ✅ Ensure TotalPrice is saved
	order.Status = "Created"  // ✅ Set Status

	createdOrder, err := s.store.CreateOrder(order)
	if err != nil {
		return models.Order{}, err
	}

	log.Println("✅ Order successfully created:", createdOrder)
	return createdOrder, nil
}


// GetOrder retrieves an order
func (s *OrderService) GetOrder(ctx context.Context, id int) (models.Order, error) {
	select {
	case <-ctx.Done():
		return models.Order{}, ctx.Err()
	default:
	}
	return s.store.GetOrder(id)
}

// UpdateOrder modifies an existing order and updates stock
func (s *OrderService) UpdateOrder(ctx context.Context, id int, updatedOrder models.Order) (models.Order, error) {
	select {
	case <-ctx.Done():
		return models.Order{}, ctx.Err()
	default:
	}

	tx, err := repositories.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Order{}, err
	}
	defer tx.Rollback()

	// Fetch existing order
	existingOrder, err := s.store.GetOrder(id)
	if err != nil {
		return models.Order{}, err
	}

	// Restore stock for old order items
	for _, item := range existingOrder.Items {
		book, err := s.bookstore.GetBook(item.BookID)
		if err != nil {
			return models.Order{}, err
		}
		book.Stock += item.Quantity
		if _, err := s.bookstore.UpdateBook(item.BookID, book); err != nil {
			return models.Order{}, err
		}
	}

	// Update stock for new order items
	for i, item := range updatedOrder.Items {
		book, err := s.bookstore.GetBook(item.BookID)
		if err != nil {
			return models.Order{}, err
		}

		if book.Stock < item.Quantity {
			return models.Order{}, fmt.Errorf("insufficient stock for book ID %d", item.BookID)
		}

		book.Stock -= item.Quantity
		if _, err := s.bookstore.UpdateBook(item.BookID, book); err != nil {
			return models.Order{}, err
		}

		updatedOrder.Items[i].BookID = book.ID
	}

	// Update order
	updatedOrder, err = s.store.UpdateOrder(id, updatedOrder)
	if err != nil {
		return models.Order{}, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return models.Order{}, err
	}

	return updatedOrder, nil
}

// DeleteOrder removes an order
func (s *OrderService) DeleteOrder(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	tx, err := repositories.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.store.DeleteOrder(id)
	if err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ListOrders fetches all orders
func (s *OrderService) ListOrders(ctx context.Context) ([]models.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.store.ListOrders()
}

// GetOrdersInRange fetches orders within a date range
func (s *OrderService) GetOrdersInRange(ctx context.Context, from, to time.Time) ([]models.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.store.GetOrdersByDateRange(from, to)
}
