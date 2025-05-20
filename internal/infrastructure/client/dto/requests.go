package dto

// GetProductRequest represents a request to get product information
type GetProductRequest struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

// GetProductResponse represents a response with product information
type GetProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

// ErrorResponse represents an error response from the product service
type ErrorResponse struct {
	Message string `json:"message"`
}
