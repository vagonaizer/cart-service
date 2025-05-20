package dto

import (
	"errors"
)

// AddItemRequest represents a request to add an item to the cart
type AddItemRequest struct {
	Count uint16 `json:"count" validate:"required,min=1"`
}

// CartItem represents an item in the cart
type CartItem struct {
	SKU      uint32 `json:"sku"`
	Quantity uint16 `json:"quantity"`
	Price    uint32 `json:"price"`
}

// GetCartResponse represents a response with cart contents
type GetCartResponse struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

// Validate validates the request
func (r *AddItemRequest) Validate() error {
	if r.Count == 0 {
		return errors.New("count must be greater than 0")
	}
	return nil
}
