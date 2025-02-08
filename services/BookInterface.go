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
		log.Fatal("❌ ERROR: BookStore or AuthorStore is nil in BookService")
	}
	return &BookService{store: bookStore, authorStore: authorStore}
}

// CreateBook inserts a new book and ensures author exists
func (bs *BookService) CreateBook(ctx context.Context, book models.Book) (models.Book, error) {
	log.Println("🔵 Inside BookService.CreateBook()")

	// Check if authorStore is nil
	if bs.authorStore == nil {
		log.Fatal("❌ ERROR: AuthorRepository (bs.authorStore) is nil!")
	}

	select {
	case <-ctx.Done():
		log.Println("❌ Context expired")
		return models.Book{}, ctx.Err()
	default:
	}
	log.Println("🔹 Received Author:", book.Author)
	log.Println("🔹 Author First Name:", book.Author.FirstName)
	log.Println("🔹 Author Last Name:", book.Author.LastName)
	log.Println("🔹 Author Bio:", book.Author.Bio)

	log.Println("🟢 Checking if author exists...")
	if book.Author == nil {
		log.Println("🛑 book.Author is nil, creating a default author object.")
		book.Author = &models.Author{}
	}

	if book.AuthorID > 0 {
		_, err := bs.authorStore.GetAuthor(book.AuthorID)
		if err != nil {
			log.Println("❌ Author does not exist")
			return models.Book{}, fmt.Errorf("author with ID %d does not exist: %w", book.AuthorID, err)
		}
	} else {
		log.Println("🟢 Creating new author...")

		if bs.authorStore == nil {
			log.Fatal("❌ ERROR: bs.authorStore is nil before calling CreateAuthor!")
		}

		if book.Author == nil {
			log.Fatal("❌ ERROR: book.Author is nil before calling CreateAuthor!")
			book.Author = &models.Author{}
		}

		// FIX: Ensure we pass a valid Author object
		newAuthor, err := bs.authorStore.CreateAuthor(models.Author{
			FirstName: book.Author.FirstName,
			LastName:  book.Author.LastName,
			Bio:       book.Author.Bio,
		})
		if err != nil {
			log.Println("❌ Failed to create author:", err)
			return models.Book{}, fmt.Errorf("failed to create author: %w", err)
		}
		book.AuthorID = newAuthor.ID
	}

	log.Println("🟢 Inserting book into DB...")

	createdBook, err := bs.store.CreateBook(book)
	if err != nil {
		log.Println("❌ Database insert failed:", err)
		return models.Book{}, err
	}

	log.Println("✅ Book successfully created:", createdBook)
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

func (bs *BookService) SearchBooks(ctx context.Context, criteria models.SearchCriteria) ([]models.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return bs.store.SearchBooks(criteria)
}
