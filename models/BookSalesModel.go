package models

import "github.com/uptrace/bun"

type BookSales struct {
	bun.BaseModel `bun:"table:book_sales"`
	ID            int  `bun:",pk,autoincrement"`
	BookID        int  `bun:",notnull"` // Foreign key to Book
	Book          *Book `bun:"rel:belongs-to,join:book_id=id"`
	Quantity      int  `bun:",notnull"`
}
