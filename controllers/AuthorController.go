package controllers

import (
	"FinalProject/models"
	"FinalProject/services"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuthorController struct {
	service *services.AuthorService
}

func NewAuthorController(s *services.AuthorService) *AuthorController {
	return &AuthorController{service: s}
}

func (ac *AuthorController) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := ac.service.CreateAuthor(ctx, author)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (ac *AuthorController) GetAuthor(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		WriteJSONError(w, http.StatusBadRequest, "invalid path")
		return
	}

	if len(parts) == 2 || (len(parts) == 3 && parts[2] == "") {
		ac.ListAuthors(w, r)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid author ID")
		return
	}

	author, getErr := ac.service.GetAuthor(ctx, id)
	if getErr != nil {
		WriteJSONError(w, http.StatusNotFound, getErr.Error())
		return
	}
	json.NewEncoder(w).Encode(author)
}

func (ac *AuthorController) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		WriteJSONError(w, http.StatusBadRequest, "missing author ID")
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid author ID")
		return
	}

	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, updateErr := ac.service.UpdateAuthor(ctx, id, author)
	if updateErr != nil {
		WriteJSONError(w, http.StatusNotFound, updateErr.Error())
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (ac *AuthorController) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		WriteJSONError(w, http.StatusBadRequest, "missing author ID")
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid author ID")
		return
	}

	delErr := ac.service.DeleteAuthor(ctx, id)
	if delErr != nil {
		WriteJSONError(w, http.StatusNotFound, delErr.Error())
		return
	}

	log.Printf("Author with ID %d successfully deleted\n", id)
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Author successfully deleted"))
}

func (ac *AuthorController) ListAuthors(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	authors, err := ac.service.ListAuthors(ctx)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(authors)
}
