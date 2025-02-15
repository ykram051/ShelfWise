package repositories

import (
	"FinalProject/models"
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

// OrderStore interface
type OrderStore interface {
	CreateOrder(o models.Order) (models.Order, error)
	GetOrder(id int) (models.Order, error)
	UpdateOrder(id int, o models.Order) (models.Order, error)
	DeleteOrder(id int) error
	ListOrders() ([]models.Order, error)
	GetOrdersByDateRange(from, to time.Time) ([]models.Order, error)
	SearchOrdersByCustomerID(customerID int) ([]models.Order, error)
}

// PostgreSQL-backed implementation of OrderStore
type OrderRepository struct {
	db *bun.DB
}

// NewOrderRepository returns a new instance
func NewOrderRepository(db *bun.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder inserts a new order
func (r *OrderRepository) CreateOrder(order models.Order) (models.Order, error) {
	// Insert Order
	_, err := r.db.NewInsert().
		Model(&order).
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return models.Order{}, fmt.Errorf("error inserting order: %w", err)
	}

	// Insert Order Items
	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		_, err := r.db.NewInsert().
			Model(&order.Items[i]).
			Returning("*").
			Exec(context.Background())
		if err != nil {
			return models.Order{}, fmt.Errorf("error inserting order item: %w", err)
		}
	}

	return order, nil
}

// GetOrder fetches an order by ID with related customer and items
func (r *OrderRepository) GetOrder(id int) (models.Order, error) {
	var order models.Order
	err := r.db.NewSelect().
		Model(&order).
		Where("?TableAlias.id = ?", id).
		Relation("Customer").
		Relation("Items.Book").
		Relation("Items.Book.Author").
		Scan(context.Background())

	if err != nil {
		return models.Order{}, fmt.Errorf("order not found: %w", err)
	}
	return order, nil
}

// GetOrdersByDateRange fetches orders in a time range
func (r *OrderRepository) GetOrdersByDateRange(from, to time.Time) ([]models.Order, error) {
	var orders []models.Order

	err := r.db.NewSelect().
		Model(&orders).
		Where("?TableAlias.created_at BETWEEN ? AND ?", from, to).
		Relation("Customer").
		Relation("Items.Book.Author").
		Scan(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error retrieving orders: %w", err)
	}
	return orders, nil
}

// ListOrders fetches all orders with relationships
func (r *OrderRepository) ListOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.NewSelect().
		Model(&orders).
		Relation("Customer").
		Relation("Items.Book").
		Relation("Items.Book.Author").
		Scan(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error retrieving orders: %w", err)
	}

	return orders, nil
}

// UpdateOrder modifies an existing order
func (r *OrderRepository) UpdateOrder(id int, order models.Order) (models.Order, error) {
	var existingOrder models.Order
	err := r.db.NewSelect().
		Model(&existingOrder).
		Where("?TableAlias.id = ?", id).
		Relation("Items").
		Scan(context.Background())

	if err != nil {
		return models.Order{}, fmt.Errorf("order with ID %d not found", id)
	}
	// to keep the original created_at and status
	order.ID = id
	order.CreatedAt = existingOrder.CreatedAt
	order.Status = existingOrder.Status

	_, err = r.db.NewDelete().
		Model((*models.OrderItem)(nil)).
		Where("order_id = ?", id).
		Exec(context.Background())

	if err != nil {
		return models.Order{}, fmt.Errorf("error clearing previous order items: %w", err)
	}

	totalPrice := 0.0
	for i := range order.Items {
		var book models.Book
		err := r.db.NewSelect().
			Model(&book).
			Where("?TableAlias.id = ?", order.Items[i].BookID).
			Scan(context.Background())

		if err != nil {
			return models.Order{}, fmt.Errorf("book with ID %d not found", order.Items[i].BookID)
		}

		order.Items[i].Book = &book
		order.Items[i].Book.PublishedAt = book.PublishedAt

		order.Items[i].OrderID = id
		_, err = r.db.NewInsert().Model(&order.Items[i]).Exec(context.Background())
		if err != nil {
			return models.Order{}, fmt.Errorf("error inserting order item: %w", err)
		}

		totalPrice += book.Price * float64(order.Items[i].Quantity)
	}
	order.TotalPrice = totalPrice

	_, err = r.db.NewUpdate().
		Model(&order).
		Where("?TableAlias.id = ?", id).
		Returning("*").
		Exec(context.Background())

	if err != nil {
		return models.Order{}, fmt.Errorf("error updating order: %w", err)
	}

	var updatedOrder models.Order
	err = r.db.NewSelect().
		Model(&updatedOrder).
		Where("?TableAlias.id = ?", id).
		Relation("Customer").
		Relation("Items.Book.Author").
		Scan(context.Background())

	if err != nil {
		return models.Order{}, fmt.Errorf("error retrieving updated order: %w", err)
	}

	return updatedOrder, nil
}

func (r *OrderRepository) DeleteOrder(id int) error {
	var order models.Order
	err := r.db.NewSelect().
		Model(&order).
		Where("id = ?", id).
		Scan(context.Background())

	if err != nil {
		return fmt.Errorf("order with ID %d not found", id)
	}
	result, err := r.db.NewDelete().
		Model((*models.Order)(nil)).
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("order with ID %d not found", id)
	}

	return nil
}

func (r *OrderRepository) SearchOrdersByCustomerID(customerID int) ([]models.Order, error) {
	var orders []models.Order

	// ✅ Step 1: Check if customer exists
	var customer models.Customer
	err := r.db.NewSelect().
		Model(&customer).
		Where("id = ?", customerID).
		Scan(context.Background())

	if err != nil {
		return nil, fmt.Errorf("customer with ID %d not found", customerID)
	}

	// ✅ Step 2: Retrieve orders for the customer
	err = r.db.NewSelect().
		Model(&orders).
		Where("?TableAlias.customer_id = ?", customerID).
		Relation("Customer").
		Relation("Items.Book.Author").
		Scan(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error retrieving orders for customer ID %d: %w", customerID, err)
	}

	return orders, nil
}
