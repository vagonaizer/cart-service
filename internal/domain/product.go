package domain

import "errors"

// Product represents a product in the system
type Product struct {
	SKU   uint32
	Name  string
	Price uint32
}

var (
	ErrProductNotFound = errors.New("product not found")
)
