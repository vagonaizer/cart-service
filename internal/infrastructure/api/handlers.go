package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"route256/cart/internal/domain/models"
	"route256/cart/internal/domain/ports"
	"route256/cart/internal/infrastructure/api/dto"
	apiErrors "route256/cart/internal/infrastructure/api/errors"
	"route256/cart/internal/usecase/cart"
)

// Handler handles HTTP requests for the cart service
type Handler struct {
	service ports.CartService
}

// NewHandler creates a new cart service handler
func NewHandler(service ports.CartService) *Handler {
	return &Handler{
		service: service,
	}
}

// AddItem handles adding an item to the cart
func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, apiErrors.ErrInvalidUserID.Error(), apiErrors.ErrInvalidUserID.Code)
		return
	}

	// Validate user_id
	if userID <= 0 {
		http.Error(w, apiErrors.ErrInvalidUserID.Error(), apiErrors.ErrInvalidUserID.Code)
		return
	}

	skuID, err := strconv.ParseUint(r.PathValue("sku_id"), 10, 32)
	if err != nil {
		http.Error(w, apiErrors.ErrInvalidSKU.Error(), apiErrors.ErrInvalidSKU.Code)
		return
	}

	// Validate sku_id
	if skuID == 0 {
		http.Error(w, apiErrors.ErrInvalidSKU.Error(), apiErrors.ErrInvalidSKU.Code)
		return
	}

	var req dto.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.AddItem(userID, uint32(skuID), req.Count); err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusPreconditionFailed)
			return
		}
		if err.Error() == "not enough items in stock" {
			http.Error(w, err.Error(), http.StatusPreconditionFailed)
			return
		}
		if apiErr, ok := apiErrors.IsAPIError(err); ok {
			http.Error(w, apiErr.Error(), apiErr.Code)
			return
		}
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveItem handles removing an item from the cart
func (h *Handler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, apiErrors.ErrInvalidUserID.Error(), apiErrors.ErrInvalidUserID.Code)
		return
	}

	// Validate user_id
	if userID <= 0 {
		http.Error(w, apiErrors.ErrInvalidUserID.Error(), apiErrors.ErrInvalidUserID.Code)
		return
	}

	skuID, err := strconv.ParseUint(r.PathValue("sku_id"), 10, 32)
	if err != nil {
		http.Error(w, apiErrors.ErrInvalidSKU.Error(), apiErrors.ErrInvalidSKU.Code)
		return
	}

	// Validate sku_id
	if skuID == 0 {
		http.Error(w, apiErrors.ErrInvalidSKU.Error(), apiErrors.ErrInvalidSKU.Code)
		return
	}

	if err := h.service.RemoveItem(userID, uint32(skuID)); err != nil {
		if apiErr, ok := apiErrors.IsAPIError(err); ok {
			http.Error(w, apiErr.Error(), apiErr.Code)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ClearCart handles clearing the cart
func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, apiErrors.ErrInvalidUserID.Error(), apiErrors.ErrInvalidUserID.Code)
		return
	}

	// Validate user_id
	if userID <= 0 {
		http.Error(w, apiErrors.ErrInvalidUserID.Error(), apiErrors.ErrInvalidUserID.Code)
		return
	}

	if err := h.service.ClearCart(userID); err != nil {
		if apiErr, ok := apiErrors.IsAPIError(err); ok {
			http.Error(w, apiErr.Error(), apiErr.Code)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetCart handles getting the cart contents
func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, apiErrors.ErrInvalidUserID.Error(), apiErrors.ErrInvalidUserID.Code)
		return
	}

	// Validate user_id
	if userID <= 0 {
		http.Error(w, apiErrors.ErrInvalidUserID.Error(), apiErrors.ErrInvalidUserID.Code)
		return
	}

	cart, err := h.service.GetCart(userID)
	if err != nil {
		if errors.Is(err, models.ErrCartNotFound) {
			http.Error(w, "cart not found", http.StatusNotFound)
			return
		}
		if apiErr, ok := apiErrors.IsAPIError(err); ok {
			http.Error(w, apiErr.Error(), apiErr.Code)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	items := make([]dto.CartItem, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = dto.CartItem{
			SKU:      item.SKU,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	resp := dto.GetCartResponse{
		Items:      items,
		TotalPrice: cart.TotalPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// Checkout handles creating an order from the cart
func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, apiErrors.ErrInvalidUserID.Error(), apiErrors.ErrInvalidUserID.Code)
		return
	}

	// Validate user_id
	if userID <= 0 {
		http.Error(w, apiErrors.ErrInvalidUserID.Error(), apiErrors.ErrInvalidUserID.Code)
		return
	}

	orderID, err := h.service.Checkout(r.Context(), userID)
	if err != nil {
		if errors.Is(err, models.ErrCartNotFound) {
			http.Error(w, "cart not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, cart.ErrCartEmpty) {
			http.Error(w, "cart is empty", http.StatusBadRequest)
			return
		}
		if apiErr, ok := apiErrors.IsAPIError(err); ok {
			http.Error(w, apiErr.Error(), apiErr.Code)
			return
		}
		log.Printf("Checkout error for user %d: %v", userID, err)
		http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	resp := dto.CheckoutResponse{
		OrderID: orderID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
