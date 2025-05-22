package models

// ItemList is a slice of Item
type ItemList []Item

// Item represents a product in the cart
type Item struct {
	// SKU is the product identifier
	SKU uint32

	// Quantity is the number of items
	Quantity uint16

	// Price is the price of one item in cents
	Price uint32
}
