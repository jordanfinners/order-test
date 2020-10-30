package model

import (
	"github.com/google/uuid"
)

// Request is a inbound http style request for a handler
type Request struct {
	Method      string
	Body        string
	QueryParams string
}

// Response is a outbound http style response for a handler
type Response struct {
	Status int
	Body   string
}

// Pack represents packs of items that can be sent to the customer
type Pack struct {
	Quantity int `bson:"quantity" json:"quantity"`
}

// Order represents a customers order and the pack(s) that fulfil it
type Order struct {
	Items int    `bson:"items" json:"items"`
	Packs []Pack `bson:"packs" json:"packs"`
}

// OrderDocument represents a customers order and the pack(s) that fulfil it with a database id
type OrderDocument struct {
	ID    uuid.UUID `bson:"id" json:"id"`
	Items int       `bson:"items" json:"items"`
	Packs []Pack    `bson:"packs" json:"packs"`
}
