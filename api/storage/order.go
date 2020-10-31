package storage

import (
	"context"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"jordanfinners/api/model"
)

const ordersCollectionName string = "orders"

// GetOrders retrieves all orders from collection
func (c Client) GetOrders(ctx context.Context) ([]model.OrderDocument, error) {
	cursor, err := c.db.Collection(ordersCollectionName).Find(ctx, bson.D{})
	if err != nil {
		log.Printf("Failed to load document cursor %v", err)
		return nil, err
	}
	var orders []model.OrderDocument
	err = cursor.All(ctx, &orders)
	if err != nil {
		log.Printf("Failed to load documents to orders %v", err)
		return nil, err
	}
	return orders, nil
}

// SaveOrder saves an order to the collection
func (c Client) SaveOrder(ctx context.Context, order model.Order) (model.OrderDocument, error) {
	document := model.OrderDocument{
		ID:    uuid.New(),
		Items: order.Items,
		Packs: order.Packs,
	}
	_, err := c.db.Collection(ordersCollectionName).InsertOne(ctx, document)
	if err != nil {
		log.Printf("Failed to save order %v", err)
		return model.OrderDocument{}, err
	}
	return document, nil
}
