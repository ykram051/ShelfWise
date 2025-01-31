package models

type OrderItem struct {
	Book     Book `json:"book"`
	Quantity int  `json:"quantity"`
}
