package repositories

import (
	"FinalProject/models"
	"context"
	"fmt"
	"strings"

	"github.com/uptrace/bun"
)

// BookStore interface
type BookStore interface {
	CreateBook(book models.Book) (models.Book, error)
	GetBook(id int) (models.Book, error)
	UpdateBook(id int, book models.Book) (models.Book, error)
	DeleteBook(id int) error
	SearchBooks(criteria models.SearchCriteria) ([]models.Book, error)
	ListBooks() ([]models.Book, error)
}

// PostgreSQL-backed implementation of BookStore
type BookRepository struct {
	db *bun.DB
}

// NewBookRepository returns a new instance
func NewBookRepository(db *bun.DB) *BookRepository {
	return &BookRepository{db: db}
}

// CreateBook inserts a new book
func (r *BookRepository) CreateBook(book models.Book) (models.Book, error) {
	_, err := r.db.NewInsert().Model(&book).Exec(context.Background())
	if err != nil {
		return models.Book{}, fmt.Errorf("error inserting book: %w", err)
	}
	return book, nil
}

// GetBook fetches a book by ID
func (r *BookRepository) GetBook(id int) (models.Book, error) {
	var book models.Book
	err := r.db.NewSelect().Model(&book).Where("id = ?", id).Relation("Author").Scan(context.Background())
	if err != nil {
		return models.Book{}, fmt.Errorf("book not found: %w", err)
	}
	return book, nil
}

// UpdateBook modifies an existing book
func (r *BookRepository) UpdateBook(id int, book models.Book) (models.Book, error) {
	book.ID = id
	_, err := r.db.NewUpdate().Model(&book).Where("id = ?", id).Exec(context.Background())
	if err != nil {
		return models.Book{}, fmt.Errorf("error updating book: %w", err)
	}
	return book, nil
}

// DeleteBook removes a book
func (r *BookRepository) DeleteBook(id int) error {
	_, err := r.db.NewDelete().Model((*models.Book)(nil)).Where("id = ?", id).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting book: %w", err)
	}
	return nil
}

// SearchBooks filters books by criteria
func (r *BookRepository) SearchBooks(criteria models.SearchCriteria) ([]models.Book, error) {
	var books []models.Book
	query := r.db.NewSelect().Model(&books).Relation("Author")

	if criteria.Title != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(criteria.Title)+"%")
	}
	if criteria.Author != "" {
		query = query.Join("JOIN authors ON authors.id = books.author_id").
			Where("LOWER(authors.first_name || ' ' || authors.last_name) LIKE ?", "%"+strings.ToLower(criteria.Author)+"%")
	}
	if criteria.Genre != "" {
		query = query.Where("? = ANY(genres)", criteria.Genre)
	}

	err := query.Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error searching books: %w", err)
	}
	return books, nil
}

// ListBooks fetches all books
func (r *BookRepository) ListBooks() ([]models.Book, error) {
	var books []models.Book
	err := r.db.NewSelect().Model(&books).Relation("Author").Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error retrieving books: %w", err)
	}
	return books, nil
}
