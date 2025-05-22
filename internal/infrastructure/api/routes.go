package api

import "net/http"

// RegisterRoutes registers all routes for the cart service
func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	// Cart operations
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", handler.AddItem)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", handler.RemoveItem)
	mux.HandleFunc("DELETE /user/{user_id}/cart", handler.ClearCart)
	mux.HandleFunc("GET /user/{user_id}/cart", handler.GetCart)
	mux.HandleFunc("POST /user/{user_id}/checkout", handler.Checkout)
}
