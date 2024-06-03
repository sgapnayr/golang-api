package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setUpRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/orders", getOrders)
	r.GET("/orders/:id", getOrderById)
	r.POST("/orders", createOrder)
	r.PUT("/orders/:id", updateOrder)
	r.DELETE("/orders/:id", deleteOrder)
	return r
}

func resetOrders() {
	orders = []Order{
		{ID: "1", Item: "Item 1", Amount: 10},
		{ID: "2", Item: "Item 2", Amount: 20},
	}
}

func TestGetOrders(t *testing.T) {
	resetOrders()
	router := setUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/orders", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Item 1")
	assert.Contains(t, w.Body.String(), "Item 2")
}

func TestGetOrderById(t *testing.T) {
	resetOrders()
	router := setUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/orders/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Item 1")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/orders/3", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "order not found")
}

func TestCreateOrder(t *testing.T) {
	resetOrders()
	router := setUpRouter()

	newOrder := `{"id":"3","item":"Item 3","amount":30}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBufferString(newOrder))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Item 3")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/orders", nil)
	router.ServeHTTP(w, req)

	assert.Contains(t, w.Body.String(), "Item 3")
}

func TestUpdateOrder(t *testing.T) {
	resetOrders()
	router := setUpRouter()

	updatedOrder := `{"id":"1","item":"Updated Item","amount":100}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/orders/1", bytes.NewBufferString(updatedOrder))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated Item")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/orders/1", nil)
	router.ServeHTTP(w, req)

	assert.Contains(t, w.Body.String(), "Updated Item")
}

func TestDeleteOrder(t *testing.T) {
	resetOrders()
	router := setUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/orders/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "order deleted")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/orders/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "order not found")
}
