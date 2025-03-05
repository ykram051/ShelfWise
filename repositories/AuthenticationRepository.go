package repositories

import (
	"FinalProject/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/uptrace/bun"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrDuplicateEmail = errors.New("email already exists")
)

// UserRepository handles user database operations
type UserRepository struct {
	DB *bun.DB
}

// NewUserRepository initializes the repository
func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser inserts a new user into the database
func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := repo.DB.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		// Check for unique constraint violation on email
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "email") {
			log.Printf("Attempted to create user with duplicate email: %s", user.Email)
			return ErrDuplicateEmail
		}
		log.Printf("Error creating user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	log.Printf("Successfully created user with email: %s", user.Email)
	return nil
}

// GetUserByEmail fetches a user by email
func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := repo.DB.NewSelect().
		Model(&user).
		Where("email = ?", email).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return &user, nil
}
