package model

// Pack represents packs of items that can be sent to the customer
type Pack struct {
	Quantity int `bson:"quantity" json:"quantity"`
}

// Order represents a customers order and the pack(s) that fulfil it
type Order struct {
	Items int    `bson:"items" json:"items"`
	Packs []Pack `bson:"packs" json:"packs"`
}
