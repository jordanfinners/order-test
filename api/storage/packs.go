package storage

import (
	"context"
	"log"
	"sort"

	"go.mongodb.org/mongo-driver/bson"

	"jordanfinners/api/model"
)

const packsCollectionName string = "packs"

// GetPacks retrieves all packs from collection
func (c Client) GetPacks(ctx context.Context) ([]model.Pack, error) {
	cursor, err := c.db.Collection(packsCollectionName).Find(ctx, bson.D{})
	if err != nil {
		log.Printf("Failed to load document cursor %v", err)
		return nil, err
	}
	var packs []model.Pack
	err = cursor.All(ctx, &packs)
	if err != nil {
		log.Printf("Failed to load documents to packs %v", err)
		return nil, err
	}

	sort.SliceStable(packs, func(i, j int) bool {
		return packs[i].Quantity > packs[j].Quantity
	})
	return packs, nil
}

// SavePack saves an pack to the collection
func (c Client) SavePack(ctx context.Context, pack model.Pack) error {
	_, err := c.db.Collection(packsCollectionName).InsertOne(ctx, pack)
	if err != nil {
		log.Printf("Failed to save pack %v", err)
		return err
	}
	return nil
}
