package controllers

import (
	"FinalProject/models"
	"FinalProject/services"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type OrderController struct {
	service *services.OrderService
}

func NewOrderController(s *services.OrderService) *OrderController {
	return &OrderController{service: s}
}

func (oc *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var order models.Order

	// Decode request JSON
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Get authenticated user info from JWT token
	authenticatedUserID, err := strconv.Atoi(r.Header.Get("X-User-ID"))
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, "Invalid authentication")
		return
	}
	authenticatedUserRole := r.Header.Get("X-User-Role")

	// Log user ID for debugging
	log.Printf("Authenticated User ID: %d, Role: %s, Request UserID: %d\n", authenticatedUserID, authenticatedUserRole, order.UserID)

	// Restrict customers to only create orders for themselves
	if authenticatedUserRole == "customer" {
		if order.UserID != authenticatedUserID {
			WriteJSONError(w, http.StatusForbidden, "Customers can only place orders for themselves")
			return
		}
	} else if authenticatedUserRole == "admin" {
		// Ensure admin-provided user_id exists
		if order.UserID == 0 {
			WriteJSONError(w, http.StatusBadRequest, "Admin must provide a valid user_id")
			return
		}
	}

	// Allow order creation
	created, err := oc.service.CreateOrder(ctx, order)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (oc *OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Extract order ID using mux
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		WriteJSONError(w, http.StatusBadRequest, "Missing order ID")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid order ID format")
		return
	}

	// Fetch order from database
	order, err := oc.service.GetOrder(ctx, id)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, "Order not found")
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (oc *OrderController) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Extract order ID using mux
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		WriteJSONError(w, http.StatusBadRequest, "Missing order ID")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, updateErr := oc.service.UpdateOrder(ctx, id, order)
	if updateErr != nil {
		WriteJSONError(w, http.StatusNotFound, updateErr.Error())
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (oc *OrderController) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Extract order ID using mux
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		WriteJSONError(w, http.StatusBadRequest, "Missing order ID")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	delErr := oc.service.DeleteOrder(ctx, id)
	if delErr != nil {
		WriteJSONError(w, http.StatusNotFound, delErr.Error()) // ✅ Return 404 if order not found
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Order with ID %d successfully deleted", id),
	})
}

func (oc *OrderController) ListOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	authenticatedUserID, err := strconv.Atoi(r.Header.Get("X-User-ID"))
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, "Invalid authentication")
		return
	}
	authenticatedUserRole := r.Header.Get("X-User-Role")

	var orders []models.Order
	if authenticatedUserRole == "admin" {
		orders, err = oc.service.ListOrders(ctx) // ✅ Admins see all orders
	} else {
		orders, err = oc.service.SearchOrdersByCustomerID(ctx, authenticatedUserID) // ✅ Customers see only their orders
	}

	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(orders)
}

func (oc *OrderController) GetOrdersByDateRange(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		WriteJSONError(w, http.StatusBadRequest, "Missing 'from' or 'to' query parameters")
		return
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid 'from' date format, expected RFC3339 (YYYY-MM-DDTHH:MM:SSZ)")
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid 'to' date format, expected RFC3339 (YYYY-MM-DDTHH:MM:SSZ)")
		return
	}

	orders, err := oc.service.GetOrdersInRange(ctx, from, to)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (oc *OrderController) SearchOrdersByCustomerID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	customerIDStr := r.URL.Query().Get("customer_id")
	if customerIDStr == "" {
		WriteJSONError(w, http.StatusBadRequest, "Missing 'customer_id' query parameter")
		return
	}

	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid customer ID format")
		return
	}

	orders, err := oc.service.SearchOrdersByCustomerID(ctx, customerID)
	if err != nil {
		if strings.Contains(err.Error(), "customer with ID") {
			WriteJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(w).Encode(orders)
}
