package handlers

import (
	"context"
	"jordanfinners/api/model"
	"jordanfinners/api/storage"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func seedPacks(testClient storage.Client) {
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 250})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 500})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 1000})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 2000})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 5000})
}

func TestMain(m *testing.M) {
	testClient := storage.StartTestDB()
	seedPacks(testClient)
	os.Exit(m.Run())
}

func TestHandlePostOrders(t *testing.T) {
	request := Request{
		Body:        `{"items":501}`,
		QueryParams: "",
	}
	response := PostOrders(request)
	require.Equal(t, http.StatusCreated, response.Status)
	expectedBody := `{"items":501,"packs":[{"quantity":500},{"quantity":250}]}`
	require.Equal(t, expectedBody, response.Body)
}

func TestHandlePostOrdersInvalidItemsOrdered(t *testing.T) {
	request := Request{
		Body:        `{"items":Seven}`,
		QueryParams: "",
	}
	response := PostOrders(request)
	require.Equal(t, http.StatusBadRequest, response.Status)
}
