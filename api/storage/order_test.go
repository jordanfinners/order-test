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
	document, err := client.SaveOrder(context.TODO(), order)
	require.NoError(t, err)

	require.Equal(t, document.Items, order.Items)
	require.Equal(t, document.Packs, order.Packs)

	orders, err := client.GetOrders(context.TODO())

	require.NoError(t, err)
	require.Len(t, orders, 1)
	require.Equal(t, []model.OrderDocument{document}, orders)
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

	require.Equal(t, order.Items, orders[0].Items)
	require.Equal(t, order.Packs, orders[0].Packs)
	require.Equal(t, order2.Items, orders[1].Items)
	require.Equal(t, order2.Packs, orders[1].Packs)

	tearDownOrders()
}
