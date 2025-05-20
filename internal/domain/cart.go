// Package domain defines core business entities for the cart service.
package domain

import "errors"

var (
	ErrCartNotFound = errors.New("cart not found")
)

// Cart represents a shopping cart containing selected items and their total price.
type Cart struct {
	// UserID uniquely identifies the user who owns the cart
	UserID int64

	// Items is the collection of products in the cart
	Items ItemList

	// TotalPrice is the sum of all items' prices in cents
	TotalPrice uint32
}

// NewCart creates a new empty cart for the given user
func NewCart(userID int64) *Cart {
	return &Cart{
		UserID: userID,
		Items:  make(ItemList, 0),
	}
}

// AddItem adds an item to the cart or updates its quantity if it exists
func (c *Cart) AddItem(item Item) {
	for i, existingItem := range c.Items {
		if existingItem.SKU == item.SKU {
			c.Items[i].Quantity += item.Quantity
			return
		}
	}
	c.Items = append(c.Items, item)
}

// RemoveItem removes an item from the cart
func (c *Cart) RemoveItem(sku uint32) {
	for i, item := range c.Items {
		if item.SKU == sku {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return
		}
	}
}

// Clear removes all items from the cart
func (c *Cart) Clear() {
	c.Items = make(ItemList, 0)
	c.TotalPrice = 0
}

// CalculateTotalPrice calculates the total price of all items in the cart
func (c *Cart) CalculateTotalPrice() {
	total := uint32(0)
	for _, item := range c.Items {
		total += item.Price * uint32(item.Quantity)
	}
	c.TotalPrice = total
}
