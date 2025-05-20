package domain

// Item represents a product in the cart
type Item struct {
	// ID uniquely identifies the product
	ID int64

	// SKU is the product's stock keeping unit
	SKU uint32

	// Quantity is the number of items in the cart
	Quantity uint16

	// Price is the price of one item in cents
	Price uint32
}

// ItemList is a slice of items
type ItemList []Item

func (list ItemList) TotalPrice() uint32 {
	total := uint32(0)
	for _, item := range list {
		total += item.Price * uint32(item.Quantity)
	}
	return total
}

func (list ItemList) TotalCount() uint16 {
	total := uint16(0)
	for _, item := range list {
		total += item.Quantity
	}
	return total
}
