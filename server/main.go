package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Order struct {
	ID     string `json:"id"`
	Item   string `json:"item"`
	Amount int    `json:"amount"`
}

var orders = []Order{
	{ID: "1", Item: "Item 1", Amount: 10},
	{ID: "2", Item: "Item 2", Amount: 20},
	{ID: "3", Item: "Item 3", Amount: 30},
}

func getOrders(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, orders)
}

func getOrderById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range orders {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "order not found"})
}

func createOrder(c *gin.Context) {
	var newOrder Order

	if err := c.BindJSON(&newOrder); err != nil {
		return
	}

	orders = append(orders, newOrder)
	c.IndentedJSON(http.StatusCreated, newOrder)
}

func updateOrder(c *gin.Context) {
	id := c.Param("id")
	var updatedOrder Order

	if err := c.BindJSON(&updatedOrder); err != nil {
		return
	}

	for i, a := range orders {
		if a.ID == id {
			orders[i] = updatedOrder
			c.IndentedJSON(http.StatusOK, updatedOrder)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "order not found"})
}

func increaseAmount(c *gin.Context) {
	id := c.Param("id")
	var amount struct {
		Amount int `json:"amount"`
	}

	if err := c.BindJSON(&amount); err != nil {
		return
	}

	for i, a := range orders {
		if a.ID == id {
			orders[i].Amount += amount.Amount
			c.IndentedJSON(http.StatusOK, orders[i])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "order not found"})
}

func decreaseAmount(c *gin.Context) {
	id := c.Param("id")
	var amount struct {
		Amount int `json:"amount"`
	}

	if err := c.BindJSON(&amount); err != nil {
		return
	}

	for i, a := range orders {
		if a.ID == id {
			orders[i].Amount -= amount.Amount
			if orders[i].Amount < 0 {
				orders[i].Amount = 0
			}
			c.IndentedJSON(http.StatusOK, orders[i])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "order not found"})
}

func deleteOrder(c *gin.Context) {
	id := c.Param("id")

	log.Printf("Delete request received for ID: %s", id)

	for i, a := range orders {
		if a.ID == id {
			orders = append(orders[:i], orders[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "order deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "order not found"})
}

func main() {
	r := gin.Default()

	// Enable CORS with custom configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/orders", getOrders)
	r.GET("/orders/:id", getOrderById)
	r.POST("/orders", createOrder)
	r.PUT("/orders/:id", updateOrder)
	r.PUT("/orders/:id/increase", increaseAmount)
	r.PUT("/orders/:id/decrease", decreaseAmount)
	r.DELETE("/orders/:id", deleteOrder)

	r.Run(":8080")
}
