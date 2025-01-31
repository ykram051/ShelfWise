package controllers

import (
	"FinalProject/models"
	"FinalProject/services"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CustomerController struct {
	service *services.CustomerService
}

func NewCustomerController(s *services.CustomerService) *CustomerController {
	return &CustomerController{service: s}
}

func (cc *CustomerController) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := cc.service.CreateCustomer(ctx, customer)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (cc *CustomerController) GetCustomer(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		WriteJSONError(w, http.StatusBadRequest, "invalid path")
		return
	}
	if len(parts) == 2 || (len(parts) == 3 && parts[2] == "") {
		cc.ListCustomers(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid customer ID")
		return
	}

	customer, getErr := cc.service.GetCustomer(ctx, id)
	if getErr != nil {
		WriteJSONError(w, http.StatusNotFound, getErr.Error())
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func (cc *CustomerController) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		WriteJSONError(w, http.StatusBadRequest, "missing customer ID")
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid customer ID")
		return
	}

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, updateErr := cc.service.UpdateCustomer(ctx, id, customer)
	if updateErr != nil {
		WriteJSONError(w, http.StatusNotFound, updateErr.Error())
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (cc *CustomerController) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		WriteJSONError(w, http.StatusBadRequest, "missing customer ID")
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid customer ID")
		return
	}

	delErr := cc.service.DeleteCustomer(ctx, id)
	if delErr != nil {
		WriteJSONError(w, http.StatusNotFound, delErr.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cc *CustomerController) ListCustomers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	customers, err := cc.service.ListCustomers(ctx)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(customers)
}
