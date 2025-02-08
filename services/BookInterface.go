package services

import (
	"FinalProject/models"
	"FinalProject/repositories"
	"context"
	"fmt"
	"log"
)

type BookService struct {
	store       repositories.BookStore
	authorStore repositories.AuthorStore
}

func NewBookService(bookStore repositories.BookStore, authorStore repositories.AuthorStore) *BookService {
	if bookStore == nil || authorStore == nil {
		log.Fatal("ERROR: BookStore or AuthorStore is nil in BookService")
	}
	return &BookService{store: bookStore, authorStore: authorStore}
}

// CreateBook inserts a new book and ensures author exists
func (bs *BookService) CreateBook(ctx context.Context, book models.Book) (models.Book, error) {

	select {
	case <-ctx.Done():
		log.Println("Context expired")
		return models.Book{}, ctx.Err()
	default:
	}

	if book.Author == nil {
		book.Author = &models.Author{}
	}

	if book.AuthorID > 0 {
		author, err := bs.authorStore.GetAuthor(book.AuthorID)
		if err != nil {
			return models.Book{}, fmt.Errorf("author with ID %d does not exist: %w", book.AuthorID, err)
		}
		book.Author = &author
	} else {

		if book.Author.FirstName == "" || book.Author.LastName == "" {
			return models.Book{}, fmt.Errorf("author first name and last name cannot be empty")
		}

		newAuthor, err := bs.authorStore.CreateAuthor(*book.Author)
		if err != nil {
			return models.Book{}, fmt.Errorf("failed to create author: %w", err)
		}

		// Assign correct AuthorID
		book.AuthorID = newAuthor.ID
		book.Author = &newAuthor
	}

	createdBook, err := bs.store.CreateBook(book)
	if err != nil {
		return models.Book{}, err
	}

	log.Println("Book successfully created:", createdBook)
	return createdBook, nil
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

	existingBook, err := bs.store.GetBook(id)
	if err != nil {
		return models.Book{}, fmt.Errorf("book with ID %d not found", id)
	}
	if book.AuthorID > 0 && book.AuthorID != existingBook.AuthorID {
		_, err := bs.authorStore.GetAuthor(book.AuthorID)
		if err != nil {
			return models.Book{}, fmt.Errorf("author with ID %d does not exist", book.AuthorID)
		}
	}
	book.ID = id
	updatedBook, err := bs.store.UpdateBook(id, book)
	if err != nil {
		return models.Book{}, fmt.Errorf("error updating book: %w", err)
	}

	return updatedBook, nil
}

func (bs *BookService) DeleteBook(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return bs.store.DeleteBook(id)
}

func (bs *BookService) SearchBooks(ctx context.Context, criteria models.SearchCriteria) ([]models.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return bs.store.SearchBooks(criteria)
}
