package cart

import (
	"errors"

	"route256/cart/internal/domain/models"
	"route256/cart/internal/domain/ports"
)

// CartService implements ports.CartService interface
type CartService struct {
	repo           ports.CartRepository
	productService ports.ProductService
}

// NewCartService creates a new cart service
func NewCartService(repo ports.CartRepository, productService ports.ProductService) ports.CartService {
	return &CartService{
		repo:           repo,
		productService: productService,
	}
}

// AddItem adds an item to the user's cart
func (s *CartService) AddItem(userID int64, sku uint32, quantity uint16) error {
	// Get product info
	product, err := s.productService.GetProduct(sku)
	if err != nil {
		return models.ErrProductNotFound
	}

	// Get or create cart
	cart, err := s.repo.GetCart(userID)
	if err != nil {
		if errors.Is(err, models.ErrCartNotFound) {
			cart = models.NewCart(userID)
			// Create new cart
			if err := s.repo.CreateCart(cart); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Add item to cart
	cart.AddItem(models.Item{
		SKU:      product.SKU,
		Quantity: quantity,
		Price:    product.Price,
	})

	// Calculate total price
	cart.CalculateTotalPrice()

	// Save cart
	return s.repo.SaveCart(cart)
}

// RemoveItem removes an item from the user's cart
func (s *CartService) RemoveItem(userID int64, sku uint32) error {
	cart, err := s.repo.GetCart(userID)
	if err != nil {
		if errors.Is(err, models.ErrCartNotFound) {
			return nil // As per spec, return success if cart doesn't exist
		}
		return err
	}

	cart.RemoveItem(sku)
	cart.CalculateTotalPrice()

	return s.repo.SaveCart(cart)
}

// ClearCart removes all items from the user's cart
func (s *CartService) ClearCart(userID int64) error {
	cart, err := s.repo.GetCart(userID)
	if err != nil {
		if errors.Is(err, models.ErrCartNotFound) {
			return nil // As per spec, return success if cart doesn't exist
		}
		return err
	}

	cart.Clear()
	return s.repo.SaveCart(cart)
}

// GetCart returns the user's cart
func (s *CartService) GetCart(userID int64) (*models.Cart, error) {
	cart, err := s.repo.GetCart(userID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, models.ErrCartNotFound
	}

	return cart, nil
}
