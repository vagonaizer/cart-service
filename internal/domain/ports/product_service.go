package ports

import "route256/cart/internal/domain/models"

// ProductService defines the interface for product operations
type ProductService interface {
	// GetProduct retrieves product information by SKU
	GetProduct(sku uint32) (*models.Product, error)
}
