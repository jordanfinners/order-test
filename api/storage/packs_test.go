package storage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"jordanfinners/api/model"
)

var client Client

func TestMain(m *testing.M) {
	client = StartTestDB()
	os.Exit(m.Run())
}

func tearDownPacks() {
	client.db.Collection(packsCollectionName).Drop(context.Background())
}

func TestGettingEmptyPacks(t *testing.T) {
	packs, err := client.GetPacks(context.TODO())

	require.NoError(t, err)
	require.Len(t, packs, 0)
}

func TestGettingPacks(t *testing.T) {
	client.SavePack(context.TODO(), model.Pack{Quantity: 1000})
	packs, err := client.GetPacks(context.TODO())

	require.NoError(t, err)
	require.Len(t, packs, 1)
	expected := []model.Pack{
		{Quantity: 1000},
	}
	require.Equal(t, expected, packs)
	tearDownPacks()
}

func TestGettingPacksSortsDescending(t *testing.T) {
	client.SavePack(context.TODO(), model.Pack{Quantity: 1000})
	client.SavePack(context.TODO(), model.Pack{Quantity: 2000})
	packs, err := client.GetPacks(context.TODO())

	require.NoError(t, err)
	require.Len(t, packs, 2)
	expected := []model.Pack{
		{Quantity: 2000},
		{Quantity: 1000},
	}
	require.Equal(t, expected, packs)
	tearDownPacks()
}
