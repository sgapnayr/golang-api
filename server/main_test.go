package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setUpRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/orders", getOrders)
	r.GET("/orders/:id", getOrderById)
	r.POST("/orders", createOrder)
	r.PUT("/orders/:id", updateOrder)
	r.PUT("/orders/:id/increase", increaseAmount)
	r.PUT("/orders/:id/decrease", decreaseAmount)
	r.DELETE("/orders/:id", deleteOrder)
	return r
}

func resetOrders() {
	orders = []Order{
		{ID: "1", Item: "Item 1", Amount: 10},
		{ID: "2", Item: "Item 2", Amount: 20},
	}
}

func stressTest(numUsers int, router *gin.Engine, t *testing.T) bool {
	var wg sync.WaitGroup

	startTime := time.Now()

	success := true
	mu := &sync.Mutex{}

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Simulate a GET request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/orders", nil)
			router.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				mu.Lock()
				success = false
				mu.Unlock()
			}

			// Simulate a POST request
			newOrder := `{"id":"3","item":"Item 3","amount":30}`
			w = httptest.NewRecorder()
			req, _ = http.NewRequest("POST", "/orders", bytes.NewBufferString(newOrder))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			if w.Code != http.StatusCreated {
				mu.Lock()
				success = false
				mu.Unlock()
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	duration := time.Since(startTime)
	log.Printf("Stress test with %d users took %s", numUsers, duration)

	return success
}

func TestIncreaseAmount(t *testing.T) {
	router := setUpRouter()
	resetOrders()

	increasePayload := `{"amount":5}`
	req, _ := http.NewRequest("PUT", "/orders/1/increase", bytes.NewBufferString(increasePayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}

	var updatedOrder Order
	if err := json.Unmarshal(w.Body.Bytes(), &updatedOrder); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if updatedOrder.Amount != 15 {
		t.Fatalf("Expected amount %d but got %d", 15, updatedOrder.Amount)
	}
}

func TestDecreaseAmount(t *testing.T) {
	router := setUpRouter()
	resetOrders()

	decreasePayload := `{"amount":5}`
	req, _ := http.NewRequest("PUT", "/orders/1/decrease", bytes.NewBufferString(decreasePayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}

	var updatedOrder Order
	if err := json.Unmarshal(w.Body.Bytes(), &updatedOrder); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if updatedOrder.Amount != 5 {
		t.Fatalf("Expected amount %d but got %d", 5, updatedOrder.Amount)
	}
}

func TestStress(t *testing.T) {
	router := setUpRouter()

	tiers := []int{1, 10, 100, 1000, 5000, 10000}
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
