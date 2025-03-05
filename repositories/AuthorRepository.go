package repositories

import (
	"FinalProject/models"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/uptrace/bun"
)

type AuthorRepository struct {
	db *bun.DB
}

// NewAuthorRepository creates an instance
func NewAuthorRepository(db *bun.DB) *AuthorRepository {
	if db == nil {
		log.Fatal("ERROR: Database connection is nil in AuthorRepository")
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

	_, err := r.db.NewInsert().Model(&author).Returning("*").Exec(context.Background())
	if err != nil {
		log.Println("Failed to insert author:", err)
		return models.Author{}, fmt.Errorf("failed to insert author: %w", err)
	}
	log.Println("Author successfully created:", author)
	return author, nil
}

// Get Author with Row-Level Locking
func (r *AuthorRepository) GetAuthor(id int) (models.Author, error) {
	var author models.Author
	err := r.db.NewSelect().
		Model(&author).
		Where("id = ?", id).
		For("UPDATE"). // Row-Level Locking
		Scan(context.Background())
	if err != nil {
		return models.Author{}, fmt.Errorf("author not found: %w", err)
	}
	return author, nil
}

func (r *AuthorRepository) UpdateAuthor(id int, author models.Author) (models.Author, error) {
	author.ID = id

	result, err := r.db.NewUpdate().
		Model(&author).
		Where("id = ?", id).
		Returning("*").
		Exec(context.Background())

	if err != nil {
		return models.Author{}, fmt.Errorf("error updating author: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return models.Author{}, fmt.Errorf("author with ID %d not found", id)
	}

	return author, nil
}

func (r *AuthorRepository) DeleteAuthor(id int) error {
	ctx := context.Background()

	log.Println("Starting deletion process for Author ID:", id)

	// ✅ Step 1: Check if the author exists
	var authorExists bool
	err := r.db.NewSelect().
		Table("authors").
		ColumnExpr("COUNT(*) > 0").
		Where("id = ?", id).
		Scan(ctx, &authorExists)

	if err != nil {
		log.Println("Database error checking author existence:", err)
		return fmt.Errorf("error checking author existence: %w", err)
	}

	if !authorExists {
		log.Println("Author not found:", id)
		return fmt.Errorf("author with ID %d not found", id)
	}

	log.Println("Author exists. Proceeding to check for books.")

	// ✅ Step 2: Check if the author has associated books
	var bookCount int
	err = r.db.NewSelect().
		Table("books").
		ColumnExpr("COUNT(*)").
		Where("author_id = ?", id).
		Scan(ctx, &bookCount)

	if err != nil {
		log.Println("Error checking associated books:", err)
		return fmt.Errorf("error checking associated books: %w", err)
	}

	log.Println("Book count for author:", id, "=", bookCount)

	if bookCount > 0 {
		log.Println("Cannot delete author:", id, "because they have", bookCount, "associated books")
		return fmt.Errorf("cannot delete author with ID %d because they have associated books", id)
	}

	log.Println("No associated books. Proceeding with deletion.")

	// ✅ Step 3: Proceed with deletion if no books exist
	result, err := r.db.NewDelete().
		Model((*models.Author)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		log.Println("Error deleting author:", err)
		return fmt.Errorf("error deleting author: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Println("Author with ID", id, "not found during deletion.")
		return fmt.Errorf("author with ID %d not found", id)
	}

	log.Println("Author successfully deleted:", id)
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

func (r *AuthorRepository) SearchAuthors(criteria models.AuthorCriteriaModel) ([]models.Author, error) {
	var authors []models.Author
	query := r.db.NewSelect().Model(&authors)

	if criteria.FirstName != "" {
		query = query.Where("LOWER(first_name) LIKE ?", "%"+strings.ToLower(criteria.FirstName)+"%")
	}

	if criteria.LastName != "" {
		query = query.Where("LOWER(last_name) LIKE ?", "%"+strings.ToLower(criteria.LastName)+"%")
	}

	err := query.Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error searching authors: %w", err)
	}
	return authors, nil
}
