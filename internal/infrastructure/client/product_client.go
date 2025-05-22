package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"route256/cart/internal/domain/models"
	"route256/cart/internal/domain/ports"
	"route256/cart/internal/infrastructure/client/dto"
)

// ProductClient implements ports.ProductService interface
type ProductClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewProductClient creates a new product service client
func NewProductClient(baseURL string, token string, httpClient *http.Client) ports.ProductService {
	return &ProductClient{
		baseURL:    baseURL,
		token:      token,
		httpClient: httpClient,
	}
}

// GetProduct implements ports.ProductService
func (c *ProductClient) GetProduct(sku uint32) (*models.Product, error) {
	reqBody := dto.GetProductRequest{
		Token: c.token,
		SKU:   sku,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/get_product", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp dto.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("product service error: %s", errResp.Message)
	}

	var productResp dto.GetProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&productResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &models.Product{
		SKU:   sku,
		Name:  productResp.Name,
		Price: productResp.Price,
	}, nil
}
