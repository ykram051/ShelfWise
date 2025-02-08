package services

import (
	"FinalProject/models"
	"FinalProject/repositories"
	"context"
	"fmt"
)

// AuthorService now interacts with DB repository
type AuthorService struct {
	authorRepo *repositories.AuthorRepository
}

func NewAuthorService(authorRepo *repositories.AuthorRepository) *AuthorService {
	return &AuthorService{authorRepo: authorRepo}
}

// CreateAuthor inserts a new author
func (s *AuthorService) CreateAuthor(ctx context.Context, author models.Author) (models.Author, error) {
	select {
	case <-ctx.Done():
		return models.Author{}, ctx.Err()
	default:
	}

	tx, err := repositories.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Author{}, err
	}
	defer tx.Rollback()

	createdAuthor, err := s.authorRepo.CreateAuthor(author)
	if err != nil {
		return models.Author{}, fmt.Errorf("error creating author: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return models.Author{}, err
	}

	return createdAuthor, nil
}

// GetAuthor retrieves an author by ID
func (s *AuthorService) GetAuthor(ctx context.Context, id int) (models.Author, error) {
	select {
	case <-ctx.Done():
		return models.Author{}, ctx.Err()
	default:
	}
	return s.authorRepo.GetAuthor(id)
}

// UpdateAuthor modifies an existing author
func (s *AuthorService) UpdateAuthor(ctx context.Context, id int, author models.Author) (models.Author, error) {
	select {
	case <-ctx.Done():
		return models.Author{}, ctx.Err()
	default:
	}

	tx, err := repositories.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Author{}, err
	}
	defer tx.Rollback()

	updatedAuthor, err := s.authorRepo.UpdateAuthor(id, author)
	if err != nil {
		return models.Author{}, fmt.Errorf("error updating author: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return models.Author{}, err
	}

	return updatedAuthor, nil
}

// DeleteAuthor removes an author
func (s *AuthorService) DeleteAuthor(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	tx, err := repositories.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.authorRepo.DeleteAuthor(id)
	if err != nil {
		return err 
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ListAuthors retrieves all authors
func (s *AuthorService) ListAuthors(ctx context.Context) ([]models.Author, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.authorRepo.ListAuthors()
}

func (s *AuthorService) SearchAuthors(ctx context.Context, criteria models.AuthorCriteriaModel) ([]models.Author, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.authorRepo.SearchAuthors(criteria)
}
