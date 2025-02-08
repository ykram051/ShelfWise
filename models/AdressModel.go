package models

type Address struct {
	Street     string `bun:"street"`
	City       string `bun:"city"`
	State      string `bun:"state"`
	PostalCode string `bun:"postal_code"`
	Country    string `bun:"country"`
}
