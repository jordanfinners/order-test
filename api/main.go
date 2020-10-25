package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

// HandleOrders deals with incoming requests for packs
func HandleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlePostOrders(w, r)
	default:
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handlePostOrders(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	packs, err := getPacks()
	if err != nil {
		log.Printf("Error loading packs: %v", err)
		http.Error(w, "Error loading packs", http.StatusInternalServerError)
		return
	}

	ordered := r.FormValue("items")
	items, err := strconv.Atoi(ordered)

	if err != nil {
		log.Printf("Error converting items ordered to int: %v", err)
		http.Error(w, "Error converting items ordered to int", http.StatusBadRequest)
		return
	}

	order := calculateOrder(items, packs)

	response, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error serialising json response: %v", err)
		http.Error(w, "Error serialising json response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	return
}

// HandleWebsite deals with incoming requests for the website
func HandleWebsite(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		serveSite(w, r)
	default:
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func serveSite(w http.ResponseWriter, r *http.Request) {

	var ordersAPI string
	ordersAPI = os.Getenv("ORDERS_API")

	template, err := template.New("website").Parse(getSiteTemplate())
	if err != nil {
		log.Printf("Error parsing index page: %v", err)
		http.Error(w, "Error parsing index page", http.StatusInternalServerError)
		return
	}

	type templateData struct {
		ActionURL string
	}

	data := templateData{
		ActionURL: ordersAPI,
	}

	err = template.Execute(w, data)
	if err != nil {
		log.Printf("Error templating index page: %v", err)
		http.Error(w, "Error templating index page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	return
}
