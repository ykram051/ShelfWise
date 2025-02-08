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
	_, err := r.db.NewInsert().Model(&order).Exec(context.Background())
	if err != nil {
		return models.Order{}, fmt.Errorf("error inserting order: %w", err)
	}
	return order, nil
}

// GetOrder fetches an order by ID with related customer and items
func (r *OrderRepository) GetOrder(id int) (models.Order, error) {
	var order models.Order
	err := r.db.NewSelect().
		Model(&order).
		Where("id = ?", id).
		Relation("Customer"). // Include customer details
		Relation("Items").    // Include order items
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
		Relation("Customer"). // Include customer info
		Relation("Items").    // Include order items
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
