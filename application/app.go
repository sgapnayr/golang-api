package application

import (
	"context"
	"fmt"
	"net/http"
)

type App struct {
	router http.Handler
}

func New() *App {
	return &App{
		router: loadRoutes(),
	}
}

func (a *App) Start(ctx context.Context) error {
	addr := "localhost:3000"
	server := &http.Server{
		Addr:    addr,
		Handler: a.router,
	}

	fmt.Printf("Starting server on %s\n", addr)
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}
