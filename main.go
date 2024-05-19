package main

import (
	"fmt"
	"context"
	"github.com/sgapnayr/golang-api/application"
)

func main() {
	app := application.New()
	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
