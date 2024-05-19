package handler

import (
	"net/http"
	"fmt"
)

type Order struct {}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating order")
	w.WriteHeader(http.StatusCreated)
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listing orders")
	w.WriteHeader(http.StatusOK)
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting order by ID")
	w.WriteHeader(http.StatusOK)
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updating order by ID")
	w.WriteHeader(http.StatusOK)
}
	
func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Deleting order by ID")	
	w.WriteHeader(http.StatusNoContent)
}