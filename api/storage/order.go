package storage

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"jordanfinners/api/model"
)

const ordersCollectionName string = "orders"

// GetOrders retrieves all orders from collection
func (c Client) GetOrders(ctx context.Context) ([]model.Order, error) {
	cursor, err := c.db.Collection(ordersCollectionName).Find(ctx, bson.D{})
	if err != nil {
		log.Printf("Failed to load document cursor %v", err)
		return nil, err
	}
	var orders []model.Order
	err = cursor.All(ctx, &orders)
	if err != nil {
		log.Printf("Failed to load documents to orders %v", err)
		return nil, err
	}
	return orders, nil
}

// SaveOrder saves an order to the collection
func (c Client) SaveOrder(ctx context.Context, order model.Order) error {
	_, err := c.db.Collection(ordersCollectionName).InsertOne(ctx, order)
	if err != nil {
		log.Printf("Failed to save pack %v", err)
		return err
	}
	return nil
}
