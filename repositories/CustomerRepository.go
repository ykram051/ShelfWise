package repositories

import (
	"FinalProject/models"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

// CustomerStore interface
type CustomerStore interface {
	GetCustomer(id int) (models.User, error)
	UpdateCustomer(id int, c models.User) (models.User, error)
	DeleteCustomer(id int) error
	ListCustomers() ([]models.User, error)
}

// PostgreSQL-backed implementation of CustomerStore
type CustomerRepository struct {
	db *bun.DB
}

func NewCustomerRepository(db *bun.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// GetCustomer fetches a User by ID
func (r *CustomerRepository) GetCustomer(id int) (models.User, error) {
	var User models.User
	err := r.db.NewSelect().Model(&User).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return models.User{}, fmt.Errorf("User not found with ID %d", id)
	}
	return User, nil
}

// UpdateCustomer modifies an existing User
func (r *CustomerRepository) UpdateCustomer(id int, User models.User) (models.User, error) {
	// Retrieve the existing User to preserve `CreatedAt`
	var existingCustomer models.User
	err := r.db.NewSelect().Model(&existingCustomer).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return models.User{}, fmt.Errorf("User with ID %d not found", id)
	}

	User.ID = id
	User.CreatedAt = existingCustomer.CreatedAt

	_, err = r.db.NewUpdate().
		Model(&User).
		Column("name", "email",
			"street", "city", "state", "postal_code", "country").
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return models.User{}, fmt.Errorf("error updating User: %w", err)
	}

	return User, nil
}

// DeleteCustomer removes a User
func (r *CustomerRepository) DeleteCustomer(id int) error {
	var User models.User
	err := r.db.NewSelect().Model(&User).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return fmt.Errorf("User with ID %d not found", id)
	}

	result, err := r.db.NewDelete().
		Model((*models.User)(nil)).
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("error deleting User: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("User with ID %d not found", id)
	}

	return nil
}

// ListCustomers fetches all customers
func (r *CustomerRepository) ListCustomers() ([]models.User, error) {
	var customers []models.User
	err := r.db.NewSelect().Model(&customers).Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error retrieving customers: %w", err)
	}
	return customers, nil
}
