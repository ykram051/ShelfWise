package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Book struct {
	bun.BaseModel `bun:"table:books"`
	ID            int       `bun:",pk,autoincrement"`
	Title         string    `bun:",notnull"`
	AuthorID      int       `bun:",notnull"` // Foreign key to Author
	Author        *Author   `bun:"rel:belongs-to,join:author_id=id"`
	Genres        []string  `bun:",array"`
	PublishedAt   time.Time `bun:",notnull"`
	Price         float64   `bun:",notnull"`
	Stock         int       `bun:",notnull"`
}
