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
)

type BookController struct {
	service *services.BookService
}

func NewBookController(s *services.BookService) *BookController {
	return &BookController{service: s}
}

func (bc *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	created, err := bc.service.CreateBook(ctx, book)
	if err != nil {
		log.Println("Service Error:", err)
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("Book created successfully:", created)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (bc *BookController) GetBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		WriteJSONError(w, http.StatusBadRequest, "missing book ID")
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	book, err := bc.service.GetBook(ctx, id)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}
	json.NewEncoder(w).Encode(book)
}

func (bc *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		WriteJSONError(w, http.StatusBadRequest, "missing book ID")
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid book ID")
		return
	}
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, updateErr := bc.service.UpdateBook(ctx, id, book)
	if updateErr != nil {
		WriteJSONError(w, http.StatusNotFound, updateErr.Error())
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (bc *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		WriteJSONError(w, http.StatusBadRequest, "missing book ID")
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	err = bc.service.DeleteBook(ctx, id)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Book with ID %d successfully deleted", id),
	})
}

func (bc *BookController) SearchBooks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")
	genre := r.URL.Query().Get("genre")

	criteria := models.SearchCriteria{
		Title:  title,
		Author: author,
		Genre:  genre,
	}

	books, err := bc.service.SearchBooks(ctx, criteria)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(books)
}
