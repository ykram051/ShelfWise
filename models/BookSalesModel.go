package models

type BookSales struct {
	Book     Book `json:"book"`
	Quantity int  `json:"quantity_sold"`
}
