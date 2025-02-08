package services

import (
	"FinalProject/models"
	"FinalProject/repositories"
	"context"
)

type CustomerService struct {
	store repositories.CustomerStore
}

func NewCustomerService(store repositories.CustomerStore) *CustomerService {
	return &CustomerService{store: store}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error) {
	select {
	case <-ctx.Done():
		return models.Customer{}, ctx.Err()
	default:
	}
	return s.store.CreateCustomer(customer)
}

func (s *CustomerService) GetCustomer(ctx context.Context, id int) (models.Customer, error) {
	select {
	case <-ctx.Done():
		return models.Customer{}, ctx.Err()
	default:
	}
	return s.store.GetCustomer(id)
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, id int, c models.Customer) (models.Customer, error) {
	select {
	case <-ctx.Done():
		return models.Customer{}, ctx.Err()
	default:
	}
	return s.store.UpdateCustomer(id, c)
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return s.store.DeleteCustomer(id)
}

func (s *CustomerService) ListCustomers(ctx context.Context) ([]models.Customer, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.store.ListCustomers()
}

