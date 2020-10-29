package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"jordanfinners/api/orders"
	"jordanfinners/api/storage"
)

// PostOrders handles POST requests to /orders
func PostOrders(request Request) Response {
	var client storage.Client = storage.NewClient()
	log.Printf(request.Body)
	var body struct {
		Items int `json:"items"`
	}
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		log.Printf("Failed unmarshal invoke request body data: %v", err)
		return Response{
			Status: http.StatusBadRequest,
			Body:   "",
		}
	}

	if body.Items < 1 {
		log.Printf("Too few items requested: %v", body.Items)
		return Response{
			Status: http.StatusBadRequest,
			Body:   "",
		}
	}

	packs, err := client.GetPacks(context.Background())
	if err != nil {
		log.Printf("Error loading packs: %v", err)
		return Response{
			Status: http.StatusInternalServerError,
			Body:   "",
		}
	}

	if len(packs) == 0 {
		log.Printf("No packs to use")
		return Response{
			Status: http.StatusGone,
			Body:   "",
		}
	}

	order := orders.CalculateOrder(body.Items, packs)

	response, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error serialising json response: %v", err)
		return Response{
			Status: http.StatusInternalServerError,
			Body:   "",
		}
	}

	return Response{
		Status: http.StatusCreated,
		Body:   string(response),
	}
}

// GetOrders handles GET requests to /orders
func GetOrders(request Request) Response {
	var client storage.Client = storage.NewClient()

	orders, err := client.GetOrders(context.Background())
	if err != nil {
		log.Printf("Error loading orders: %v", err)
		return Response{
			Status: http.StatusInternalServerError,
			Body:   "",
		}
	}

	response, err := json.Marshal(orders)
	if err != nil {
		log.Printf("Error serialising json response: %v", err)
		return Response{
			Status: http.StatusInternalServerError,
			Body:   "",
		}
	}

	return Response{
		Status: http.StatusOK,
		Body:   string(response),
	}
}
