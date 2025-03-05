package models

import (
	"time"

	"github.com/uptrace/bun"
)

// User represents a user in the system
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID           int       `json:"id" bun:",pk,autoincrement"`
	Name         string    `json:"name" bun:",notnull"`
	Email        string    `json:"email" bun:",unique,notnull"`
	PasswordHash string    `json:"-" bun:",notnull"` // Hide from JSON
	Role         string    `json:"role" bun:",notnull"`
	Address      Address   `json:"address" bun:",embed"` 
	CreatedAt    time.Time `json:"created_at" bun:",default:current_timestamp"`
}
