package repositories

import (
	"FinalProject/models"
	"context"
	"fmt"
	"log"
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

	log.Println("✅ Order inserted with ID:", order.ID)

	// Insert Order Items
	for i := range order.Items {
		order.Items[i].OrderID = order.ID // ✅ Assign correct OrderID
		_, err := r.db.NewInsert().
			Model(&order.Items[i]).
			Returning("*"). // ✅ Return the inserted row with ID
			Exec(context.Background())
		if err != nil {
			return models.Order{}, fmt.Errorf("error inserting order item: %w", err)
		}
		log.Println("✅ Inserted order item:", order.Items[i])
	}

	return order, nil
}

// GetOrder fetches an order by ID with related customer and items
func (r *OrderRepository) GetOrder(id int) (models.Order, error) {
	var order models.Order
	err := r.db.NewSelect().
		Model(&order).
		Where("?TableAlias.id = ?", id). // ✅ Uses Bun's table alias instead of hardcoding "orders"
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
		Where("created_at BETWEEN ? AND ?", from, to).
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
		Relation("Customer").          // ✅ Fetch customer info
		Relation("Items.Book").        // ✅ Fetch books in the order items
		Relation("Items.Book.Author"). // ✅ Fetch book authors
		Scan(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error retrieving orders: %w", err)
	}

	return orders, nil
}

// UpdateOrder modifies an existing order
func (r *OrderRepository) UpdateOrder(id int, order models.Order) (models.Order, error) {
	_, err := r.db.NewUpdate().
		Model(&order).
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return models.Order{}, fmt.Errorf("error updating order: %w", err)
	}
	return order, nil
}

func (r *OrderRepository) DeleteOrder(id int) error {
	_, err := r.db.NewDelete().Model((*models.Order)(nil)).Where("id = ?", id).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}
	return nil
}
