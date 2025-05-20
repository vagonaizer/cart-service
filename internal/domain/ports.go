package domain

// CartRepository defines the interface for cart storage
type CartRepository interface {
	CreateCart(cart *Cart) error
	GetCart(userID int64) (*Cart, error)
	SaveCart(cart *Cart) error
	DeleteCart(userID int64) error
}

// ProductService defines the interface for product information
type ProductService interface {
	GetProduct(sku uint32) (*Product, error)
}

// CartService defines the interface for cart business logic
type CartService interface {
	AddItem(userID int64, sku uint32, quantity uint16) error
	RemoveItem(userID int64, sku uint32) error
	ClearCart(userID int64) error
	GetCart(userID int64) (*Cart, error)
}
