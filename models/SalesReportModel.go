package models

import (
	"time"

	"github.com/uptrace/bun"
)

type SalesReport struct {
	bun.BaseModel    `bun:"table:sales_reports"`
	ID               int         `bun:",pk,autoincrement"`
	Timestamp        time.Time   `bun:",nullzero,notnull,default:current_timestamp"`
	TotalRevenue     float64     `bun:",notnull"`
	TotalOrders      int         `bun:",notnull"`
	TopSellingBooks  []BookSales `bun:"rel:has-many,join:id=book_id"`
}
