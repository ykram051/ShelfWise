package services

import (
	"FinalProject/models"
	"FinalProject/repositories"
	"context"
	"fmt"
	"log"
)

type AuthorService struct {
	authorStore repositories.AuthorStore
	bookStore   repositories.BookStore
}

func NewAuthorService(authorStore repositories.AuthorStore, bookStore repositories.BookStore) *AuthorService {
	return &AuthorService{authorStore: authorStore, bookStore: bookStore}
}

func (s *AuthorService) CreateAuthor(ctx context.Context, author models.Author) (models.Author, error) {
	select {
	case <-ctx.Done():
		return models.Author{}, ctx.Err()
	default:
	}
	return s.authorStore.CreateAuthor(author)
}

func (s *AuthorService) GetAuthor(ctx context.Context, id int) (models.Author, error) {
	select {
	case <-ctx.Done():
		return models.Author{}, ctx.Err()
	default:
	}
	return s.authorStore.GetAuthor(id)
}

func (s *AuthorService) UpdateAuthor(ctx context.Context, id int, author models.Author) (models.Author, error) {
	select {
	case <-ctx.Done():
		return models.Author{}, ctx.Err()
	default:
	}
	return s.authorStore.UpdateAuthor(id, author)
}

func (s *AuthorService) DeleteAuthor(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		log.Printf("Error: Context canceled while deleting author ID %d\n", id)
		return ctx.Err()
	default:
		// Check the existence of the author 
		_, err := s.authorStore.GetAuthor(id)
		if err != nil {
			log.Printf("Error: Author with ID %d not found: %v\n", id, err)
			return fmt.Errorf("author not found: %v", err)
		}

		books, err := s.bookStore.ListBooks()
		if err != nil {
			log.Printf("Error: Failed to list books while deleting author ID %d: %v\n", id, err)
			return fmt.Errorf("failed to list books: %v", err)
		}
		for _, book := range books {
			if book.Author.ID == id {
				log.Printf("Error: Cannot delete author ID %d because books are associated with this author\n", id)
				return fmt.Errorf("cannot delete author with ID %d: books are associated with this author", id)
			}
		}

		if err := s.authorStore.DeleteAuthor(id); err != nil {
			log.Printf("Error: Failed to delete author ID %d: %v\n", id, err)
			return fmt.Errorf("failed to delete author: %v", err)
		}

		log.Printf("Author with ID %d successfully deleted\n", id)
		return nil
	}
}

func (s *AuthorService) ListAuthors(ctx context.Context) ([]models.Author, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.authorStore.ListAuthors()
}

func (s *AuthorService) SaveChanges() error {
	return s.authorStore.Save()
}
