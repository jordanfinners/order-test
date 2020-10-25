package main

import (
	"context"
	"log"
	"os"

	"github.com/jordanfinners/api"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {

	port := "8080"

	err := os.Setenv("ORDERS_API", "http://localhost:8080/orders")
	if err != nil {
		log.Printf("Failed to set the ORDERS_API: %v", err)
	}

	ctx := context.Background()
	err = funcframework.RegisterHTTPFunctionContext(ctx, "/orders", api.HandleOrders)
	if err != nil {
		log.Printf("funcframework.RegisterHTTPFunctionContext orders: %v", err)
	}

	err = funcframework.RegisterHTTPFunctionContext(ctx, "/website", api.HandleWebsite)
	if err != nil {
		log.Printf("funcframework.RegisterHTTPFunctionContext website: %v", err)
	}

	err = funcframework.Start(port)
	if err != nil {
		log.Printf("funcframework.Start: %v", err)
	}
}
