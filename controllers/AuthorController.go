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

	author.ID = id

	updatedAuthor, updateErr := ac.service.UpdateAuthor(ctx, id, author)
	if updateErr != nil {
		WriteJSONError(w, http.StatusNotFound, updateErr.Error())
		return
	}

	json.NewEncoder(w).Encode(updatedAuthor)
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

	err = ac.service.DeleteAuthor(ctx, id)
	if err != nil {
		WriteJSONError(w, http.StatusConflict, "error deleting author because the author has associated books")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Author with ID %d successfully deleted", id),
	})
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

func (ac *AuthorController) SearchAuthors(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	criteria := models.AuthorCriteriaModel{
		FirstName: firstName,
		LastName:  lastName,
	}

	authors, err := ac.service.SearchAuthors(ctx, criteria)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}
