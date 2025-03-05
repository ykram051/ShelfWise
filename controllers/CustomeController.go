package controllers

import (
	"FinalProject/models"
	"FinalProject/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type CustomerController struct {
	service *services.CustomerService
}

func NewCustomerController(s *services.CustomerService) *CustomerController {
	return &CustomerController{service: s}
}



func (cc *CustomerController) GetCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Extract customer ID using mux.Vars()
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		WriteJSONError(w, http.StatusBadRequest, "Missing User ID")
		return
	}

	// Convert customer ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid User ID format")
		return
	}

	// Fetch customer from the service
	customer, getErr := cc.service.GetCustomer(ctx, id)
	if getErr != nil {
		WriteJSONError(w, http.StatusNotFound, getErr.Error())
		return
	}

	// Return customer data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (cc *CustomerController) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		WriteJSONError(w, http.StatusBadRequest, "missing User ID")
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid User ID")
		return
	}

	var User models.User
	if err := json.NewDecoder(r.Body).Decode(&User); err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, updateErr := cc.service.UpdateCustomer(ctx, id, User)
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
		WriteJSONError(w, http.StatusBadRequest, "missing User ID")
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid User ID")
		return
	}

	err = cc.service.DeleteCustomer(ctx, id)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("User with ID %d successfully deleted", id),
	})
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
