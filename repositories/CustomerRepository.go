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
		return models.Customer{}, fmt.Errorf("customer not found with ID %d", id)
	}
	return customer, nil
}

// UpdateCustomer modifies an existing customer
func (r *CustomerRepository) UpdateCustomer(id int, customer models.Customer) (models.Customer, error) {
	// Retrieve the existing customer to preserve `CreatedAt`
	var existingCustomer models.Customer
	err := r.db.NewSelect().Model(&existingCustomer).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return models.Customer{}, fmt.Errorf("customer with ID %d not found", id)
	}

	customer.ID = id
	customer.CreatedAt = existingCustomer.CreatedAt

	_, err = r.db.NewUpdate().
		Model(&customer).
		Column("name", "email",
			"street", "city", "state", "postal_code", "country").
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return models.Customer{}, fmt.Errorf("error updating customer: %w", err)
	}

	return customer, nil
}

// DeleteCustomer removes a customer
func (r *CustomerRepository) DeleteCustomer(id int) error {
	var customer models.Customer
	err := r.db.NewSelect().Model(&customer).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return fmt.Errorf("customer with ID %d not found", id)
	}

	result, err := r.db.NewDelete().
		Model((*models.Customer)(nil)).
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("error deleting customer: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("customer with ID %d not found", id)
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
