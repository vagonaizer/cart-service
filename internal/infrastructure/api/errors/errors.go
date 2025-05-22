package errors

import (
	"errors"
	"net/http"
)

// APIError represents an API error
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

var (
	ErrInvalidUserID = &APIError{
		Code:    http.StatusBadRequest,
		Message: "invalid user_id",
	}

	ErrInvalidSKU = &APIError{
		Code:    http.StatusBadRequest,
		Message: "invalid sku_id",
	}

	ErrCartNotFound = &APIError{
		Code:    http.StatusNotFound,
		Message: "cart not found",
	}

	ErrItemNotFound = &APIError{
		Code:    http.StatusNotFound,
		Message: "item not found",
	}
)

// IsAPIError checks if an error is an API error
func IsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}
