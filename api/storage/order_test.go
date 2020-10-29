package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"jordanfinners/api/model"
)

func tearDownOrders() {
	client.db.Collection(ordersCollectionName).Drop(context.Background())
}

func TestGettingEmptyOrders(t *testing.T) {
	packs, err := client.GetOrders(context.TODO())

	require.NoError(t, err)
	require.Len(t, packs, 0)
}

func TestSavingOrder(t *testing.T) {
	order := model.Order{
		Items: 251,
		Packs: []model.Pack{
			{Quantity: 500},
		},
	}
	client.SaveOrder(context.TODO(), order)
	orders, err := client.GetOrders(context.TODO())

	require.NoError(t, err)
	require.Len(t, orders, 1)
	require.Equal(t, []model.Order{order}, orders)
	tearDownOrders()
}

func TestGettingOrders(t *testing.T) {
	order := model.Order{
		Items: 251,
		Packs: []model.Pack{
			{Quantity: 500},
		},
	}
	order2 := model.Order{
		Items: 10,
		Packs: []model.Pack{
			{Quantity: 251},
		},
	}
	client.SaveOrder(context.TODO(), order)
	client.SaveOrder(context.TODO(), order2)
	orders, err := client.GetOrders(context.TODO())

	require.NoError(t, err)
	require.Len(t, orders, 2)
	expected := []model.Order{
		order,
		order2,
	}
	require.Equal(t, expected, orders)
	tearDownOrders()
}
