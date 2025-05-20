package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"route256/cart/config"
	"route256/cart/internal/core"
	"route256/cart/internal/infrastructure/api"
)

func main() {
	// Parse command line arguments
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create application
	app := core.NewApp(cfg)

	// Create HTTP server
	server := &http.Server{
		Addr:              ":8082",
		Handler:           api.LoggingMiddleware(app.Mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start server
	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
