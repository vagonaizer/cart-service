package inmemory

import (
	"sync"

	"route256/cart/internal/domain/models"
)

// CartRepository implements domain.CartRepository interface
// using in-memory storage
type CartRepository struct {
	mu    sync.RWMutex
	carts map[int64]*models.Cart
}

// NewCartRepository creates a new in-memory cart repository
func NewCartRepository() *CartRepository {
	return &CartRepository{
		carts: make(map[int64]*models.Cart),
	}
}

// CreateCart implements domain.CartRepository
func (r *CartRepository) CreateCart(cart *models.Cart) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.carts[cart.UserID]; exists {
		return models.ErrCartAlreadyExists
	}

	r.carts[cart.UserID] = cart
	return nil
}

// GetCart implements domain.CartRepository
func (r *CartRepository) GetCart(userID int64) (*models.Cart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cart, exists := r.carts[userID]
	if !exists {
		return nil, models.ErrCartNotFound
	}

	return cart, nil
}

// SaveCart implements domain.CartRepository
func (r *CartRepository) SaveCart(cart *models.Cart) error {
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
		return models.ErrCartNotFound
	}

	delete(r.carts, userID)
	return nil
}
