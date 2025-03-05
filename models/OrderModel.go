package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders"`
	ID            int         `bun:",pk,autoincrement"`
	UserID        int         `bun:",notnull"` // Foreign key to User (replaces CustomerID)
	User          *User       `bun:"rel:belongs-to,join:user_id=id"` // Relationship to User
	Items         []OrderItem `bun:"rel:has-many,join:id=order_id"`  // Relationship to OrderItem
	TotalPrice    float64     `bun:",notnull"`
	CreatedAt     time.Time   `bun:",nullzero,notnull,default:current_timestamp"`
	Status        string      `bun:",notnull"`
}
