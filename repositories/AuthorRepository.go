package repositories

import (
	"FinalProject/models"
	"context"
	"fmt"
	"log"

	"github.com/uptrace/bun"
)

type AuthorRepository struct {
	db *bun.DB
}

// NewAuthorRepository creates an instance
func NewAuthorRepository(db *bun.DB) *AuthorRepository {
	if db == nil {
		log.Fatal("❌ ERROR: Database connection is nil in AuthorRepository")
	}
	return &AuthorRepository{db: db}
}

type AuthorStore interface {
	CreateAuthor(author models.Author) (models.Author, error)
	GetAuthor(id int) (models.Author, error)
	UpdateAuthor(id int, author models.Author) (models.Author, error)
	DeleteAuthor(id int) error
	ListAuthors() ([]models.Author, error)
}

func (r *AuthorRepository) CreateAuthor(author models.Author) (models.Author, error) {
	// Ensure first_name & last_name are being passed
	log.Println("🔹 Inserting Author:", author.FirstName, author.LastName, author.Bio)

	
	_, err := r.db.NewInsert().Model(&author).Returning("*").Exec(context.Background())
	if err != nil {
		log.Println("❌ Failed to insert author:", err)
		return models.Author{}, fmt.Errorf("failed to insert author: %w", err)
	}
	log.Println("✅ Author successfully created:", author)
	return author, nil
}

// Get Author with Row-Level Locking
func (r *AuthorRepository) GetAuthor(id int) (models.Author, error) {
	var author models.Author
	err := r.db.NewSelect().
		Model(&author).
		Where("id = ?", id).
		For("UPDATE"). // Prevents concurrent updates on the same row
		Scan(context.Background())
	if err != nil {
		return models.Author{}, fmt.Errorf("author not found: %w", err)
	}
	return author, nil
}

func (r *AuthorRepository) UpdateAuthor(id int, author models.Author) (models.Author, error) {
	_, err := r.db.NewUpdate().
		Model(&author).
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return models.Author{}, fmt.Errorf("error updating author: %w", err)
	}
	return author, nil
}
func (r *AuthorRepository) DeleteAuthor(id int) error {
	_, err := r.db.NewDelete().
		Model((*models.Author)(nil)).
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting author: %w", err)
	}
	return nil
}

func (r *AuthorRepository) ListAuthors() ([]models.Author, error) {
	var authors []models.Author
	err := r.db.NewSelect().Model(&authors).Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error retrieving authors: %w", err)
	}
	return authors, nil
}
