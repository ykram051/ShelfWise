package models

import "github.com/uptrace/bun"

type OrderItem struct {
	bun.BaseModel `bun:"table:order_items"`
	ID            int  `bun:",pk,autoincrement"`
	OrderID       int  `bun:",notnull"` // Foreign key to Order
	BookID        int  `bun:",notnull"` // Foreign key to Book
	Book          *Book `bun:"rel:belongs-to,join:book_id=id"`
	Quantity      int  `bun:",notnull"`
}
