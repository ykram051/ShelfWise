package services

import (
	"FinalProject/models"
	"FinalProject/repositories"
	"context"
	"fmt"
)

type BookService struct {
	store       repositories.BookStore
	authorStore repositories.AuthorStore
}

func NewBookService(s repositories.BookStore, authorStore repositories.AuthorStore) *BookService {
	return &BookService{store: s,
		authorStore: authorStore}

}

func (bs *BookService) CreateBook(ctx context.Context, book models.Book) (models.Book, error) {
	select {
	case <-ctx.Done():
		return models.Book{}, ctx.Err()
	default:

		if book.Author.ID > 0 {
			_, err := bs.authorStore.GetAuthor(book.Author.ID)
			if err != nil {
				return models.Book{}, fmt.Errorf("author with ID %d does not exist: %w", book.Author.ID, err)
			}
		} else {
			newAuthor, err := bs.authorStore.CreateAuthor(book.Author)
			if err != nil {
				return models.Book{}, fmt.Errorf("failed to create author: %w", err)
			}
			book.Author = newAuthor
		}

		createdBook, err := bs.store.CreateBook(book)
		if err != nil {
			return models.Book{}, err
		}

		return createdBook, nil
	}
}

func (bs *BookService) GetBook(ctx context.Context, id int) (models.Book, error) {
	select {
	case <-ctx.Done():
		return models.Book{}, ctx.Err()
	default:
	}
	return bs.store.GetBook(id)
}

func (bs *BookService) UpdateBook(ctx context.Context, id int, book models.Book) (models.Book, error) {
	select {
	case <-ctx.Done():
		return models.Book{}, ctx.Err()
	default:
	}
	return bs.store.UpdateBook(id, book)
}

func (bs *BookService) DeleteBook(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return bs.store.DeleteBook(id)
}

func (bs *BookService) SearchBooks(ctx context.Context, criteria repositories.SearchCriteria) ([]models.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return bs.store.SearchBooks(criteria)
}

func (bs *BookService) SaveChanges() error {
	return bs.store.Save()
}
