package app

import (
	"net/http"
	"time"

	"route256/cart/config"
	"route256/cart/internal/domain/ports"
	"route256/cart/internal/infrastructure/api"
	"route256/cart/internal/infrastructure/client"
	"route256/cart/internal/infrastructure/loms"
	"route256/cart/internal/infrastructure/repository/inmemory"
	"route256/cart/internal/usecase/cart"
)

// App represents the application
type App struct {
	Mux     *http.ServeMux
	Service ports.CartService
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config) *App {
	// Create HTTP client with retry middleware
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: client.NewRetryMiddleware(
			http.DefaultTransport,
			3,
			time.Second,
		),
	}

	// Create product service client
	productClient := client.NewProductClient(
		cfg.ProductService.URL,
		cfg.ProductService.Token,
		httpClient,
	)

	// Create LOMS client
	lomsClient, err := loms.NewClient(cfg.LOMS.Address)
	if err != nil {
		panic(err)
	}

	// Create in-memory cart repository
	repo := inmemory.NewCartRepository()

	// Create cart service
	cartService := cart.NewCartService(repo, productClient, lomsClient)

	// Create HTTP router
	mux := http.NewServeMux()
	handler := api.NewHandler(cartService)
	api.RegisterRoutes(mux, handler)

	return &App{
		Mux:     mux,
		Service: cartService,
	}
}
