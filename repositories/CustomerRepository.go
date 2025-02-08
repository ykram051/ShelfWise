package repositories

import (
	"FinalProject/models"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

// CustomerStore interface
type CustomerStore interface {
	CreateCustomer(c models.Customer) (models.Customer, error)
	GetCustomer(id int) (models.Customer, error)
	UpdateCustomer(id int, c models.Customer) (models.Customer, error)
	DeleteCustomer(id int) error
	ListCustomers() ([]models.Customer, error)
}

// PostgreSQL-backed implementation of CustomerStore
type CustomerRepository struct {
	db *bun.DB
}

// NewCustomerRepository returns a new instance
func NewCustomerRepository(db *bun.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// CreateCustomer inserts a new customer
func (r *CustomerRepository) CreateCustomer(customer models.Customer) (models.Customer, error) {
	_, err := r.db.NewInsert().Model(&customer).Exec(context.Background())
	if err != nil {
		return models.Customer{}, fmt.Errorf("error inserting customer: %w", err)
	}
	return customer, nil
}

// GetCustomer fetches a customer by ID
func (r *CustomerRepository) GetCustomer(id int) (models.Customer, error) {
	var customer models.Customer
	err := r.db.NewSelect().Model(&customer).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return models.Customer{}, fmt.Errorf("customer not found: %w", err)
	}
	return customer, nil
}

// UpdateCustomer modifies an existing customer
func (r *CustomerRepository) UpdateCustomer(id int, customer models.Customer) (models.Customer, error) {
	customer.ID = id
	_, err := r.db.NewUpdate().Model(&customer).Where("id = ?", id).Exec(context.Background())
	if err != nil {
		return models.Customer{}, fmt.Errorf("error updating customer: %w", err)
	}
	return customer, nil
}

// DeleteCustomer removes a customer
func (r *CustomerRepository) DeleteCustomer(id int) error {
	_, err := r.db.NewDelete().Model((*models.Customer)(nil)).Where("id = ?", id).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting customer: %w", err)
	}
	return nil
}

// ListCustomers fetches all customers
func (r *CustomerRepository) ListCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	err := r.db.NewSelect().Model(&customers).Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error retrieving customers: %w", err)
	}
	return customers, nil
}
