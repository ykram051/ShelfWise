package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders"`
	ID            int         `bun:",pk,autoincrement"`
	CustomerID    int         `bun:",notnull"` // Foreign key to Customer
	Customer      *Customer   `bun:"rel:belongs-to,join:customer_id=id"`
	Items         []OrderItem `bun:"rel:has-many,join:id=order_id"`
	TotalPrice    float64     `bun:",notnull"`
	CreatedAt     time.Time   `bun:",nullzero,notnull,default:current_timestamp"`
	Status        string      `bun:",notnull"`
}
