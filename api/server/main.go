package main

import (
	"log"
	"net/http"
	"os"

	"jordanfinners/api/router"
)

func main() {
	customHandlerPort, _ := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("/orders", router.HandleOrdersRequest)
	log.Fatal(http.ListenAndServe(":"+customHandlerPort, mux))
}
