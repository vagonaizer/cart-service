package ports

import "route256/cart/internal/domain/models"

// CartRepository defines the interface for cart storage operations
type CartRepository interface {
	// GetCart retrieves a cart by user ID
	GetCart(userID int64) (*models.Cart, error)

	// SaveCart saves or updates a cart
	SaveCart(cart *models.Cart) error

	// CreateCart creates a new empty cart
	CreateCart(cart *models.Cart) error
}
