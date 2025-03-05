package models

type Address struct {
	Street     string `json:"street" bun:"street"`
	City       string `json:"city" bun:"city"`
	State      string `json:"state" bun:"state"`
	PostalCode string `json:"postal_code" bun:"postal_code"`
	Country    string `json:"country" bun:"country"`
}
