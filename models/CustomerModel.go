package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Customer struct {
	bun.BaseModel `bun:"table:customers"`
	ID            int       `bun:",pk,autoincrement"`
	Name          string    `bun:",notnull"`
	Email         string    `bun:",unique,notnull"`
	Address       Address   `bun:",embed"` // Embedded Address
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
