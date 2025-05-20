package api

import (
	"net/http"
	"strings"
)

// Router handles routing for the cart service
type Router struct {
	handler *Handler
}

// NewRouter creates a new router
func NewRouter(handler *Handler) *Router {
	return &Router{
		handler: handler,
	}
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := strings.Trim(req.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 3 || parts[0] != "user" {
		http.NotFound(w, req)
		return
	}

	userID := parts[1]
	if userID == "" {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	req.URL.RawQuery = "user_id=" + userID + "&" + req.URL.RawQuery

	switch {
	case len(parts) == 3 && parts[2] == "cart":
		switch req.Method {
		case http.MethodGet:
			r.handler.GetCart(w, req)
		case http.MethodDelete:
			r.handler.ClearCart(w, req)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	case len(parts) == 4 && parts[2] == "cart":
		req.URL.RawQuery = "sku_id=" + parts[3] + "&" + req.URL.RawQuery
		switch req.Method {
		case http.MethodPost:
			r.handler.AddItem(w, req)
		case http.MethodDelete:
			r.handler.RemoveItem(w, req)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.NotFound(w, req)
	}
}
