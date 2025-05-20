package ports

import "route256/cart/internal/domain/models"

// CartService defines the interface for cart operations
type CartService interface {
	// AddItem adds an item to the cart
	AddItem(userID int64, sku uint32, count uint16) error

	// RemoveItem removes an item from the cart
	RemoveItem(userID int64, sku uint32) error

	// GetCart retrieves the cart contents
	GetCart(userID int64) (*models.Cart, error)

	// ClearCart removes all items from the cart
	ClearCart(userID int64) error
}
