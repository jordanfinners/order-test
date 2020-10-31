package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"jordanfinners/api/model"
	"jordanfinners/api/orders"
	"jordanfinners/api/storage"
)

// PostOrders handles POST requests to /orders
func PostOrders(request model.Request) model.Response {
	var client storage.Client = storage.NewClient()

	var body struct {
		Items int `json:"items"`
	}
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		log.Printf("Failed unmarshal invoke request body data: %v", err)
		return model.Response{
			Status: http.StatusBadRequest,
			Body:   "",
		}
	}

	if body.Items < 1 {
		log.Printf("Too few items requested: %v", body.Items)
		return model.Response{
			Status: http.StatusBadRequest,
			Body:   "",
		}
	}

	packs, err := client.GetPacks(context.Background())
	if err != nil {
		log.Printf("Error loading packs: %v", err)
		return model.Response{
			Status: http.StatusInternalServerError,
			Body:   "",
		}
	}

	if len(packs) == 0 {
		log.Printf("No packs to use")
		return model.Response{
			Status: http.StatusGone,
			Body:   "",
		}
	}

	order := orders.CalculateOrder(body.Items, packs)

	document, err := client.SaveOrder(context.Background(), order)
	if err != nil {
		log.Printf("Error saving order: %v", err)
		return model.Response{
			Status: http.StatusInternalServerError,
			Body:   "",
		}
	}

	response, err := json.Marshal(document)
	if err != nil {
		log.Printf("Error serialising json response: %v", err)
		return model.Response{
			Status: http.StatusInternalServerError,
			Body:   "",
		}
	}

	return model.Response{
		Status: http.StatusCreated,
		Body:   string(response),
	}
}

// GetOrders handles GET requests to /orders
func GetOrders(request model.Request) model.Response {
	var client storage.Client = storage.NewClient()

	orders, err := client.GetOrders(context.Background())
	if err != nil {
		log.Printf("Error loading orders: %v", err)
		return model.Response{
			Status: http.StatusInternalServerError,
			Body:   "",
		}
	}

	response, err := json.Marshal(orders)
	if err != nil {
		log.Printf("Error serialising json response: %v", err)
		return model.Response{
			Status: http.StatusInternalServerError,
			Body:   "",
		}
	}

	return model.Response{
		Status: http.StatusOK,
		Body:   string(response),
	}
}
