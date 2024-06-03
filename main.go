package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
)

type Order struct {
	ID     string `json:"id"`
	Item   string `json:"item"`
	Amount int    `json:"amount"`
}

var orders = []Order{
	{ID: "1", Item: "Item 1", Amount: 10},
	{ID: "2", Item: "Item 2", Amount: 20},
}

var (
	orderCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "orders_total",
			Help: "Total number of orders processed",
		},
		[]string{"method"},
	)
	kafkaWriter *kafka.Writer
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[*websocket.Conn]bool)
	broadcast = make(chan Order)
)

func init() {
	prometheus.MustRegister(orderCounter)
	initKafka()
	go handleMessages()
}

func initKafka() {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}
}

func produceMessage(order Order) {
	message := kafka.Message{
		Key:   []byte(order.ID),
		Value: []byte(fmt.Sprintf("Order: %s, Item: %s, Amount: %d", order.ID, order.Item, order.Amount)),
	}
	err := kafkaWriter.WriteMessages(context.Background(), message)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
	}
}

func handleMessages() {
	for {
		order := <-broadcast
		for client := range clients {
			err := client.WriteJSON(order)
			if err != nil {
				log.Printf("WebSocket error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func wsEndpoint(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		var order Order
		err := conn.ReadJSON(&order)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			delete(clients, conn)
			break
		}
	}
}

func getOrders(c *gin.Context) {
	orderCounter.WithLabelValues("GET").Inc()
	c.IndentedJSON(http.StatusOK, orders)
}

func getOrderById(c *gin.Context) {
	orderCounter.WithLabelValues("GET").Inc()
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
	produceMessage(newOrder)
	orderCounter.WithLabelValues("POST").Inc()
	broadcast <- newOrder
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
			produceMessage(updatedOrder)
			orderCounter.WithLabelValues("PUT").Inc()
			broadcast <- updatedOrder
			c.IndentedJSON(http.StatusOK, updatedOrder)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "order not found"})
}

func deleteOrder(c *gin.Context) {
	id := c.Param("id")

	for i, a := range orders {
		if a.ID == id {
			orders = append(orders[:i], orders[i+1:]...)
			deletedOrder := Order{ID: id, Item: "deleted"}
			produceMessage(deletedOrder)
			orderCounter.WithLabelValues("DELETE").Inc()
			broadcast <- deletedOrder
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
	r.DELETE("/orders/:id", deleteOrder)
	r.GET("/ws", wsEndpoint)

	// Expose Prometheus metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(":8080")
}
