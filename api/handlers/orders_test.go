package handlers

import (
	"context"
	"encoding/json"
	"jordanfinners/api/model"
	"jordanfinners/api/storage"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var testClient storage.Client

func seedPacks(testClient storage.Client) {
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 250})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 500})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 1000})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 2000})
	testClient.SavePack(context.TODO(), model.Pack{Quantity: 5000})
}

func TestMain(m *testing.M) {
	testClient = storage.StartTestDB()
	seedPacks(testClient)
	os.Exit(m.Run())
}

func TestHandleGetOrders(t *testing.T) {
	order := model.Order{
		Items: 10,
		Packs: []model.Pack{
			{Quantity: 251},
		},
	}
	testClient.SaveOrder(context.TODO(), order)
	request := model.Request{
		Body:        "",
		QueryParams: "",
	}
	response := GetOrders(request)
	require.Equal(t, http.StatusOK, response.Status)

	var body []model.OrderDocument
	err := json.Unmarshal([]byte(response.Body), &body)
	require.NoError(t, err)

	require.Equal(t, order.Items, body[0].Items)
	require.Equal(t, order.Packs, body[0].Packs)
}

func TestHandlePostOrders(t *testing.T) {
	request := model.Request{
		Body:        `{"items":250}`,
		QueryParams: "",
	}
	response := PostOrders(request)
	require.Equal(t, http.StatusCreated, response.Status)

	var body model.OrderDocument
	err := json.Unmarshal([]byte(response.Body), &body)
	require.NoError(t, err)

	require.Equal(t, 250, body.Items)

	expectedPacks := []model.Pack{{Quantity: 250}}
	require.Equal(t, expectedPacks, body.Packs)
}

func TestHandlePostOrdersInvalidItemsOrdered(t *testing.T) {
	request := model.Request{
		Body:        `{"items":Seven}`,
		QueryParams: "",
	}
	response := PostOrders(request)
	require.Equal(t, http.StatusBadRequest, response.Status)
}

func TestHandlePostOrdersZeroItemsOrdered(t *testing.T) {
	request := model.Request{
		Body:        `{"items":0}`,
		QueryParams: "",
	}
	response := PostOrders(request)
	require.Equal(t, http.StatusBadRequest, response.Status)
}
