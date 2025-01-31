package models

import "time"

type SalesReport struct {
	Timestamp       time.Time   `json:"timestamp"`
	TotalRevenue    float64     `json:"total_revenue"`
	TotalOrders     int         `json:"total_orders"`
	TopSellingBooks []BookSales `json:"top_selling_books"`
}
