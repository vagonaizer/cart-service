package cart

import (
	"context"
	"errors"

	"route256/cart/internal/domain/models"
	"route256/cart/internal/domain/ports"
)

var (
	ErrCartEmpty = errors.New("cart is empty")
)

// CartService implements ports.CartService interface
type CartService struct {
	repo           ports.CartRepository
	productService ports.ProductService
	lomsClient     ports.LOMSClient
}

// NewCartService creates a new cart service
func NewCartService(repo ports.CartRepository, productService ports.ProductService, lomsClient ports.LOMSClient) ports.CartService {
	return &CartService{
		repo:           repo,
		productService: productService,
		lomsClient:     lomsClient,
	}
}

// AddItem adds an item to the user's cart
func (s *CartService) AddItem(userID int64, sku uint32, quantity uint16) error {
	// Get product info
	product, err := s.productService.GetProduct(sku)
	if err != nil {
		return models.ErrProductNotFound
	}

	// Check stock quantity
	stock, err := s.lomsClient.GetStocksInfo(context.Background(), sku)
	if err != nil {
		return err
	}

	// Get current cart to check existing items
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

	// Calculate total quantity including existing items
	totalQuantity := quantity
	for _, item := range cart.Items {
		if item.SKU == sku {
			totalQuantity += item.Quantity
		}
	}

	if uint64(totalQuantity) > stock {
		return errors.New("not enough items in stock")
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

// Checkout creates an order from the cart and clears it
func (s *CartService) Checkout(ctx context.Context, userID int64) (int64, error) {
	// Get cart
	cart, err := s.repo.GetCart(userID)
	if err != nil {
		if errors.Is(err, models.ErrCartNotFound) {
			return 0, models.ErrCartNotFound
		}
		return 0, err
	}

	// Check if cart is empty
	if len(cart.Items) == 0 {
		return 0, ErrCartEmpty
	}

	// Convert cart items to LOMS items
	items := make([]ports.Item, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = ports.Item{
			SKU:   item.SKU,
			Count: item.Quantity, // Quantity is already uint16
		}
	}

	// Create order in LOMS
	orderID, err := s.lomsClient.CreateOrder(ctx, userID, items)
	if err != nil {
		return 0, err
	}

	// Check order status
	orderInfo, err := s.lomsClient.GetOrderInfo(ctx, orderID)
	if err != nil {
		return 0, err
	}

	if orderInfo.Status == "failed" {
		return 0, errors.New("order creation failed: not enough items in stock")
	}

	// Clear cart after successful order creation
	if err := s.ClearCart(userID); err != nil {
		return 0, err
	}

	return orderID, nil
}
