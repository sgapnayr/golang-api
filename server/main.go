package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Order struct {
	ID     string `json:"id" dynamodbav:"id"`
	Item   string `json:"item" dynamodbav:"item"`
	Amount int    `json:"amount" dynamodbav:"amount"`
}

var svc *dynamodb.DynamoDB

func main() {
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

	// Start the server
	r.Run(":8080")
}

// Handlers
func createOrderHandler(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a new UUID for the order ID
	order.ID = uuid.New().String()

	existingOrder, _ := getOrder(svc, order.ID, order.Item)
	if existingOrder != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order with this ID and Item already exists"})
		return
	}

	err := createOrder(svc, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order_id": order.ID})
}

func getOrderHandler(c *gin.Context) {
	id := c.Param("id")
	item := c.Param("item")
	order, err := getOrder(svc, id, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func getAllOrdersHandler(c *gin.Context) {
	orders, err := getAllOrders(svc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func updateOrderHandler(c *gin.Context) {
	id := c.Param("id")
	item := c.Param("item")
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.ID = id   // Ensure the ID in the path matches the ID in the JSON body
	order.Item = item // Ensure the Item in the path matches the Item in the JSON body
	err := updateOrder(svc, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

func deleteOrderHandler(c *gin.Context) {
	id := c.Param("id")
	item := c.Param("item")
	err := deleteOrder(svc, id, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// DynamoDB CRUD operations
func createOrder(svc *dynamodb.DynamoDB, order Order) error {
	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		return fmt.Errorf("got error marshalling new order item: %v", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("orders"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return fmt.Errorf("got error calling PutItem: %v", err)
	}
	return nil
}

func getOrder(svc *dynamodb.DynamoDB, id string, item string) (*Order, error) {
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("orders"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"item": {
				S: aws.String(item),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("got error calling GetItem: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	order := new(Order)
	err = dynamodbattribute.UnmarshalMap(result.Item, order)
	if err != nil {
		return nil, fmt.Errorf("got error unmarshalling order: %v", err)
	}

	return order, nil
}

func getAllOrders(svc *dynamodb.DynamoDB) ([]Order, error) {
	// Use the DynamoDB Scan API to get all items in the table
	input := &dynamodb.ScanInput{
		TableName: aws.String("orders"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("got error calling Scan: %v", err)
	}

	var orders []Order
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &orders)
	if err != nil {
		return nil, fmt.Errorf("got error unmarshalling orders: %v", err)
	}

	return orders, nil
}

func updateOrder(svc *dynamodb.DynamoDB, order Order) error {
	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		return fmt.Errorf("got error marshalling order item: %v", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("orders"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return fmt.Errorf("got error calling PutItem: %v", err)
	}
	return nil
}

func deleteOrder(svc *dynamodb.DynamoDB, id string, item string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("orders"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"item": {
				S: aws.String(item),
			},
		},
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("got error calling DeleteItem: %v", err)
	}
	return nil
}
