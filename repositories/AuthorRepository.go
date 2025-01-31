package repositories

import (
	"FinalProject/models"
	"errors"
	"fmt"
	"sync"
)

type AuthorStore interface {
	CreateAuthor(author models.Author) (models.Author, error)
	GetAuthor(id int) (models.Author, error)
	UpdateAuthor(id int, author models.Author) (models.Author, error)
	DeleteAuthor(id int) error
	ListAuthors() ([]models.Author, error)
	Save() error
}

type InMemoryAuthorStore struct {
	authors map[int]models.Author
	nextID  int
	mu      sync.Mutex
	backend string // "authors.json"
}

func NewInMemoryAuthorStore(fileName string) *InMemoryAuthorStore {
	store := &InMemoryAuthorStore{
		authors: make(map[int]models.Author),
		backend: fileName,
	}
	LoadFromFile(fileName, &store.authors, &store.mu)

	for id := range store.authors {
		if id > store.nextID {
			store.nextID = id
		}
	}
	return store
}

func (s *InMemoryAuthorStore) CreateAuthor(author models.Author) (models.Author, error) {
	s.mu.Lock()
	s.nextID++
	author.ID = s.nextID
	s.authors[author.ID] = author
	s.mu.Unlock()

	if err := s.Save(); err != nil {
		fmt.Printf("Failed to save author data: %v\n", err)
	}
	return author, nil
}

func (s *InMemoryAuthorStore) GetAuthor(id int) (models.Author, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.authors[id]
	if !ok {
		return models.Author{}, errors.New("author not found")
	}
	return a, nil
}

func (s *InMemoryAuthorStore) UpdateAuthor(id int, author models.Author) (models.Author, error) {
	s.mu.Lock()
	if _, ok := s.authors[id]; !ok {
		return models.Author{}, errors.New("author not found")
	}
	author.ID = id
	s.authors[id] = author
	s.mu.Unlock()
	if err := s.Save(); err != nil {
		return models.Author{}, fmt.Errorf("failed to save after updating author: %v", err)
	}

	return author, nil
}

func (s *InMemoryAuthorStore) DeleteAuthor(id int) error {
	s.mu.Lock()

	if _, ok := s.authors[id]; !ok {
		return errors.New("author not found")
	}
	delete(s.authors, id)
	s.mu.Unlock()
	if err := s.Save(); err != nil {
		return fmt.Errorf("failed to save after deleting author: %v", err)
	}
	return nil
}

func (s *InMemoryAuthorStore) ListAuthors() ([]models.Author, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []models.Author
	for _, a := range s.authors {
		result = append(result, a)
	}
	return result, nil
}

func (s *InMemoryAuthorStore) Save() error {
	if err := SaveToFile(s.backend, s.authors, &s.mu); err != nil {
		return fmt.Errorf("error saving authors: %v", err)
	}
	return nil
}
