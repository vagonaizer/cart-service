package models

// Product represents a product in the catalog
type Product struct {
	// SKU is the product identifier
	SKU uint32

	// Name is the product name
	Name string

	// Price is the product price in cents
	Price uint32
}
