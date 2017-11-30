package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products", productsHandler).Methods("GET")
	r.HandleFunc("/product/{id:[0-9]+}", productByIDHandler).Methods("GET")
	r.HandleFunc("/product/{id:[0-9]+}", deleteProductHandler).Methods("DELETE")
	r.HandleFunc("/product", createOrUpdateProductHander).Methods("PUT")

	server := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func productByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	p, _, err := findProduct(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func createOrUpdateProductHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var p Product
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newP, err := createOrUpdateProduct(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(newP)

}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	deleted, err := deleteProduct(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(deleted)

}

func deleteProduct(id string) (Product, error) {
	toDelete, idx, err := findProduct(id)
	if err != nil {
		var empty Product
		return empty, err
	}
	deleted := *toDelete //Because of slice pointers (defer the pointer before array change)
	data = append(data[:idx], data[idx+1:]...)
	return deleted, nil
}

// Stored data (repository) functions
func createOrUpdateProduct(p Product) (Product, error) {
	if p.ID == "" {
		p.ID = strconv.Itoa(len(data) + 1)
		p.CreatedAt = time.Now()
		data = append(data, p)
		return p, nil
	}
	stored, _, err := findProduct(p.ID)
	if err != nil {
		var empty Product
		return empty, err
	}

	//C hange found product
	stored.ID = p.ID
	stored.Name = p.Name
	stored.Price = p.Price
	return *stored, nil
}

func findProduct(id string) (*Product, int, error) {
	for idx := range data {
		p := &data[idx]
		if p.ID == id {
			return p, idx, nil
		}
	}
	return nil, -1, fmt.Errorf("Product with id %s does not exist", id)
}

// Model

// Product model record
type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Price     float64   `json:"price"`
}

// Category model record
type Category struct {
	id   string
	name string
}

var data = []Product{
	Product{ID: "1", Name: "Product 1", CreatedAt: time.Now(), Price: 100.0},
	Product{ID: "2", Name: "Product 2", CreatedAt: time.Now(), Price: 100.0},
	Product{ID: "3", Name: "Product 3", CreatedAt: time.Now(), Price: 100.0},
	Product{ID: "4", Name: "Product 4", CreatedAt: time.Now(), Price: 100.0},
	Product{ID: "5", Name: "Product 5", CreatedAt: time.Now(), Price: 100.0},
}
