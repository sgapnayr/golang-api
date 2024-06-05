package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrderHandler(t *testing.T) {
	// Mock the request payload
	order := Order{
		Item:   "Item1",
		Amount: 10,
	}
	payload, _ := json.Marshal(order)

	// Create a request
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// Initialize the Gin engine
	r := setupRouter()

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Order created successfully")

	// Extract the generated order ID from the response
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["order_id"])
}

func TestGetOrderHandler(t *testing.T) {
	// First, create an order to ensure there is an order to get
	order := Order{
		Item:   "Item1",
		Amount: 10,
	}
	payload, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	r := setupRouter()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	orderID := response["order_id"].(string)

	// Mock the request to get the created order
	req, _ = http.NewRequest("GET", fmt.Sprintf("/orders/%s/Item1", orderID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAllOrdersHandler(t *testing.T) {
	// Mock the request
	req, _ := http.NewRequest("GET", "/orders", nil)

	// Initialize the Gin engine
	r := setupRouter()

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateOrderHandler(t *testing.T) {
	// First, create an order to ensure there is an order to update
	order := Order{
		Item:   "Item1",
		Amount: 10,
	}
	payload, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	r := setupRouter()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	orderID := response["order_id"].(string)

	// Mock the request payload for update
	order.Amount = 20
	payload, _ = json.Marshal(order)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/orders/%s/Item1", orderID), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Order updated successfully")
}

func TestDeleteOrderHandler(t *testing.T) {
	// First, create an order to ensure there is an order to delete
	order := Order{
		Item:   "Item1",
		Amount: 10,
	}
	payload, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	r := setupRouter()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	orderID := response["order_id"].(string)

	// Mock the request to delete the created order
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/orders/%s/Item1", orderID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Order deleted successfully")
}

func TestStress(t *testing.T) {
	router := setupRouter()

	tiers := []int{1, 10, 25}
	for _, numUsers := range tiers {
		t.Run(fmt.Sprintf("%d users", numUsers), func(t *testing.T) {
			resetOrders()
			success := stressTest(numUsers, router, t)
			if success {
				fmt.Printf("Test with %d users: PASS\n", numUsers)
			} else {
				fmt.Printf("Test with %d users: FAIL\n", numUsers)
			}
		})
	}
}

func stressTest(numUsers int, router *gin.Engine, t *testing.T) bool {
	var wg sync.WaitGroup
	success := true
	mu := sync.Mutex{}

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			order := Order{
				Item:   fmt.Sprintf("item-%d", i),
				Amount: i,
			}
			payload, _ := json.Marshal(order)
			req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			mu.Lock()
			if w.Code != http.StatusOK {
				success = false
				t.Errorf("Failed to create order %d: %v", i, w.Body.String())
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return success
}

func resetOrders() {
	// Implement a function to reset the orders in your database if needed.
	// This could be deleting all orders or resetting the DynamoDB table.
}

func setupRouter() *gin.Engine {
	// Initialize AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))

	// Initialize DynamoDB client
	svc = dynamodb.New(sess)

	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.Default())

	// Define routes
	r.POST("/orders", createOrderHandler)
	r.GET("/orders/:id/:item", getOrderHandler)
	r.GET("/orders", getAllOrdersHandler)
	r.PUT("/orders/:id/:item", updateOrderHandler)
	r.DELETE("/orders/:id/:item", deleteOrderHandler)

	return r
}
