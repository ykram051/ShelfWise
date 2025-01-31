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
		WriteJSONError(w, http.StatusNotFound, delErr.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
