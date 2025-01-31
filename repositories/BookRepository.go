package repositories

import (
	"FinalProject/models"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type SearchCriteria struct {
	Title  string
	Author string
	Genre  string
}

type BookStore interface {
	CreateBook(book models.Book) (models.Book, error)
	GetBook(id int) (models.Book, error)
	UpdateBook(id int, book models.Book) (models.Book, error)
	DeleteBook(id int) error
	SearchBooks(criteria SearchCriteria) ([]models.Book, error)
	ListBooks() ([]models.Book, error)
	Save() error
}

type InMemoryBookStore struct {
	books   map[int]models.Book
	nextID  int
	mu      sync.Mutex
	backend string // "books.json"
}

func NewInMemoryBookStore(fileName string) *InMemoryBookStore {
	store := &InMemoryBookStore{
		books:   make(map[int]models.Book),
		backend: fileName,
	}
	LoadFromFile(fileName, &store.books, &store.mu)

	for id := range store.books {
		if id > store.nextID {
			store.nextID = id
		}
	}
	return store
}

func (s *InMemoryBookStore) CreateBook(book models.Book) (models.Book, error) {
	s.mu.Lock()
	s.nextID++
	book.ID = s.nextID
	s.books[book.ID] = book
	s.mu.Unlock()
	err := s.Save()
	if err != nil {
		return models.Book{}, fmt.Errorf("error saving books: %v", err)
	}
	return book, nil
}

func (s *InMemoryBookStore) GetBook(id int) (models.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book, ok := s.books[id]
	if !ok {
		return models.Book{}, errors.New("book not found")
	}

	return book, nil
}

func (s *InMemoryBookStore) UpdateBook(id int, book models.Book) (models.Book, error) {
	s.mu.Lock()
	if _, ok := s.books[id]; !ok {
		s.mu.Unlock()
		return models.Book{}, errors.New("book not found")
	}
	book.ID = id
	s.books[id] = book
	s.mu.Unlock()

	err := s.Save()
	if err != nil {
		return models.Book{}, fmt.Errorf("error saving the updated books: %v", err)
	}
	return book, nil
}

func (s *InMemoryBookStore) DeleteBook(id int) error {
	s.mu.Lock()

	if _, ok := s.books[id]; !ok {
		s.mu.Unlock()
		return errors.New("book not found")
	}
	delete(s.books, id)
	s.mu.Unlock()
	err := s.Save()
	if err != nil {
		return fmt.Errorf("error saving books after the delete: %v", err)
	}
	return nil
}

func (s *InMemoryBookStore) SearchBooks(criteria SearchCriteria) ([]models.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var results []models.Book
	for _, book := range s.books {
		if criteria.Title != "" && !strings.Contains(strings.ToLower(book.Title), strings.ToLower(criteria.Title)) {
			continue
		}
		fullName := strings.ToLower(book.Author.FirstName + " " + book.Author.LastName)
		if criteria.Author != "" && !strings.Contains(fullName, strings.ToLower(criteria.Author)) {
			continue
		}
		if criteria.Genre != "" {
			match := false
			for _, g := range book.Genres {
				if strings.EqualFold(g, criteria.Genre) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		results = append(results, book)
	}
	return results, nil
}
func (s *InMemoryBookStore) ListBooks() ([]models.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []models.Book
	for _, b := range s.books {
		result = append(result, b)
	}
	return result, nil
}

func (s *InMemoryBookStore) Save() error {
	if err := SaveToFile(s.backend, s.books, &s.mu); err != nil {
		return fmt.Errorf("error saving books: %v", err)
	}
	return nil
}
