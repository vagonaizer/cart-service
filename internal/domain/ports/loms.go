package ports

import (
	"context"
)

// LOMSClient defines the interface for interacting with the LOMS service
type LOMSClient interface {
	// CreateOrder creates a new order from cart items
	CreateOrder(ctx context.Context, userID int64, items []Item) (int64, error)

	// GetStocksInfo checks if there are enough items in stock
	GetStocksInfo(ctx context.Context, sku uint32) (uint64, error)

	// GetOrderInfo retrieves information about an order
	GetOrderInfo(ctx context.Context, orderID int64) (*OrderInfo, error)
}

// Item represents a cart item for order creation
type Item struct {
	SKU   uint32
	Count uint16
}

// OrderInfo represents information about an order
type OrderInfo struct {
	Status string
	UserID int64
	Items  []Item
}
