package models


import "github.com/uptrace/bun"

type Author struct {
	bun.BaseModel `bun:"table:authors"`
	ID        int     `bun:",pk,autoincrement"`
	FirstName     string `bun:",notnull"`
	LastName      string `bun:",notnull"`
	Bio           string
}
