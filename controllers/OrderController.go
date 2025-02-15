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
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := oc.service.CreateOrder(ctx, order)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (oc *OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		WriteJSONError(w, http.StatusBadRequest, "invalid path")
		return
	}

	if len(parts) == 2 || (len(parts) == 3 && parts[2] == "") {
		oc.ListOrders(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid order ID")
		return
	}

	order, getErr := oc.service.GetOrder(ctx, id)
	if getErr != nil {
		WriteJSONError(w, http.StatusNotFound, getErr.Error())
		return
	}
	json.NewEncoder(w).Encode(order)
}

func (oc *OrderController) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		WriteJSONError(w, http.StatusBadRequest, "missing order ID")
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid order ID")
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

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		WriteJSONError(w, http.StatusBadRequest, "missing order ID")
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid order ID")
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

	orders, err := oc.service.ListOrders(ctx)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(orders)
}

func (oc *OrderController) GetOrdersByDateRange(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// ✅ Parse "from" and "to" query parameters
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	// ✅ Validate input dates
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

	// ✅ Fetch orders within the given date range
	orders, err := oc.service.GetOrdersInRange(ctx, from, to)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// ✅ Return the filtered orders
	json.NewEncoder(w).Encode(orders)
}

func (oc *OrderController) SearchOrdersByCustomerID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// ✅ Parse `customer_id` from query parameters
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

	// ✅ Fetch orders for the given Customer ID
	orders, err := oc.service.SearchOrdersByCustomerID(ctx, customerID)
	if err != nil {
		if strings.Contains(err.Error(), "customer with ID") { // ✅ Detect customer not found error
			WriteJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// ✅ Return the filtered orders
	json.NewEncoder(w).Encode(orders)
}
