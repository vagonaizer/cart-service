# Cart Service

## Overview
Cart Service is a microservice responsible for managing shopping carts in an e-commerce system. It provides functionality for adding items to cart, removing items, clearing cart, and checkout operations.

## Architecture

### Clean Architecture
The service follows Clean Architecture principles with clear separation of concerns:

- **Domain Layer**: Contains core business logic and interfaces
  - Models: `Cart`, `Item`, `Product`
  - Ports: Interfaces for repositories and external services

- **Application Layer**: Implements use cases
  - Cart service implementation
  - Business logic orchestration

- **Infrastructure Layer**: Implements external interfaces
  - HTTP API handlers
  - In-memory repository
  - External service clients (Product Service, LOMS)

### Key Components

#### HTTP API
- RESTful API implementation using standard `net/http`
- Middleware for logging and error handling
- Input validation and error responses
- Graceful shutdown implementation

Example of graceful shutdown from `main.go`:
```go
// Channel to listen for an interrupt or terminate signal from the OS.
shutdown := make(chan os.Signal, 1)
signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

// Give outstanding requests 5 seconds to complete.
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

#### Repository
- In-memory implementation with thread-safe operations
- Mutex-based concurrency control
- CRUD operations for cart management

Example of thread-safe repository:
```go
type CartRepository struct {
    mu    sync.RWMutex
    carts map[int64]*models.Cart
}
```

#### External Service Integration
- Product Service client with retry mechanism
- LOMS (Logistics and Order Management System) integration
- HTTP client with timeout and retry middleware

Example of retry middleware:
```go
httpClient := &http.Client{
    Timeout: 5 * time.Second,
    Transport: client.NewRetryMiddleware(
        http.DefaultTransport,
        3,
        time.Second,
    ),
}
```

## API Endpoints

### Cart Management
- `POST /api/v1/cart/{user_id}/items/{sku_id}` - Add item to cart
- `DELETE /api/v1/cart/{user_id}/items/{sku_id}` - Remove item from cart
- `DELETE /api/v1/cart/{user_id}` - Clear cart
- `GET /api/v1/cart/{user_id}` - Get cart contents
- `POST /api/v1/cart/{user_id}/checkout` - Checkout cart

## Error Handling
- Custom error types for different scenarios
- HTTP status codes mapping
- Detailed error messages
- Input validation

## Configuration
- YAML-based configuration
- Environment-specific settings
- Service endpoints configuration
- Timeout and retry settings

## Technologies Used
- Go 1.x
- Standard library `net/http` for HTTP server
- Protocol Buffers for service communication
- In-memory storage with mutex synchronization
- Retry mechanism for external service calls

## Best Practices Implemented
1. **Clean Architecture**
   - Clear separation of concerns
   - Dependency injection
   - Interface-based design

2. **Error Handling**
   - Custom error types
   - Proper error wrapping
   - HTTP status code mapping

3. **Concurrency**
   - Thread-safe operations
   - Mutex-based synchronization
   - Graceful shutdown

4. **API Design**
   - RESTful endpoints
   - Input validation
   - Consistent error responses

5. **Reliability**
   - Retry mechanism for external calls
   - Timeout handling
   - Graceful degradation

## Running the Service
```bash
go run cmd/cart/main.go -config config/config.yaml
```

## Testing
- Unit tests for repository
- Integration tests for handlers
- Benchmark tests for performance