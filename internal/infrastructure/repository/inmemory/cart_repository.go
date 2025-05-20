package inmemory

import (
	"errors"
	"sync"

	"route256/cart/internal/domain"
)

// CartRepository implements domain.CartRepository interface
// using in-memory storage
type CartRepository struct {
	mu    sync.RWMutex
	carts map[int64]*domain.Cart
}

// NewCartRepository creates a new in-memory cart repository
func NewCartRepository() *CartRepository {
	return &CartRepository{
		carts: make(map[int64]*domain.Cart),
	}
}

// CreateCart implements domain.CartRepository
func (r *CartRepository) CreateCart(cart *domain.Cart) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.carts[cart.UserID]; exists {
		return errors.New("cart already exists")
	}

	r.carts[cart.UserID] = cart
	return nil
}

// GetCart implements domain.CartRepository
func (r *CartRepository) GetCart(userID int64) (*domain.Cart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cart, exists := r.carts[userID]
	if !exists {
		return nil, domain.ErrCartNotFound
	}

	return cart, nil
}

// SaveCart implements domain.CartRepository
func (r *CartRepository) SaveCart(cart *domain.Cart) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.carts[cart.UserID] = cart
	return nil
}

// DeleteCart implements domain.CartRepository
func (r *CartRepository) DeleteCart(userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.carts[userID]; !exists {
		return domain.ErrCartNotFound
	}

	delete(r.carts, userID)
	return nil
}
