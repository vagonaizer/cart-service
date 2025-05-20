package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"route256/cart/internal/domain/models"
	"route256/cart/internal/infrastructure/api"
	"route256/cart/internal/infrastructure/repository/inmemory"
	"route256/cart/internal/usecase/cart"
)

func setupTestServer(t *testing.T) (*httptest.Server, *api.Handler) {
	repo := inmemory.NewCartRepository()
	// В интеграционных тестах используем реальный ProductClient
	productService := &mockProductService{
		products: map[uint32]*models.Product{
			123: {
				SKU:   123,
				Name:  "Test Product",
				Price: 1000,
			},
		},
	}
	cartService := cart.NewCartService(repo, productService)
	handler := api.NewHandler(cartService)

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, handler)

	server := httptest.NewServer(mux)
	return server, handler
}

// mockProductService для интеграционных тестов
type mockProductService struct {
	products map[uint32]*models.Product
}

func (m *mockProductService) GetProduct(sku uint32) (*models.Product, error) {
	if product, ok := m.products[sku]; ok {
		return product, nil
	}
	return nil, models.ErrProductNotFound
}

func TestCartAPI_AddItem(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.Close()

	tests := []struct {
		name       string
		userID     int64
		sku        uint32
		count      uint16
		wantStatus int
	}{
		{
			name:       "success",
			userID:     1,
			sku:        123,
			count:      2,
			wantStatus: http.StatusOK,
		},
		{
			name:       "product not found",
			userID:     1,
			sku:        999,
			count:      1,
			wantStatus: http.StatusPreconditionFailed,
		},
		{
			name:       "invalid user id",
			userID:     0,
			sku:        123,
			count:      1,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := map[string]uint16{
				"count": tt.count,
			}
			jsonBody, err := json.Marshal(reqBody)
			require.NoError(t, err)

			url := server.URL + "/user/" + strconv.FormatInt(tt.userID, 10) + "/cart/" + strconv.FormatUint(uint64(tt.sku), 10)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestCartAPI_GetCart(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.Close()

	// Сначала добавим товар в корзину
	userID := int64(1)
	sku := uint32(123)
	count := uint16(2)

	// Add item
	reqBody := map[string]uint16{
		"count": count,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	url := server.URL + "/user/" + strconv.FormatInt(userID, 10) + "/cart/" + strconv.FormatUint(uint64(sku), 10)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()

	// Get cart
	url = server.URL + "/user/" + strconv.FormatInt(userID, 10) + "/cart"
	req, err = http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var cartResp struct {
		Items []struct {
			SKU      uint32 `json:"sku"`
			Quantity uint16 `json:"quantity"`
			Price    uint32 `json:"price"`
		} `json:"items"`
		TotalPrice uint32 `json:"total_price"`
	}

	err = json.NewDecoder(resp.Body).Decode(&cartResp)
	require.NoError(t, err)

	assert.Len(t, cartResp.Items, 1)
	assert.Equal(t, sku, cartResp.Items[0].SKU)
	assert.Equal(t, count, cartResp.Items[0].Quantity)
	assert.Equal(t, uint32(1000), cartResp.Items[0].Price)
	assert.Equal(t, uint32(2000), cartResp.TotalPrice)
}

func TestCartAPI_RemoveItem(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.Close()

	// Сначала добавим товар в корзину
	userID := int64(1)
	sku := uint32(123)
	count := uint16(2)

	// Add item
	reqBody := map[string]uint16{
		"count": count,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	url := server.URL + "/user/" + strconv.FormatInt(userID, 10) + "/cart/" + strconv.FormatUint(uint64(sku), 10)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()

	// Remove item
	url = server.URL + "/user/" + strconv.FormatInt(userID, 10) + "/cart/" + strconv.FormatUint(uint64(sku), 10)
	req, err = http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify cart is empty
	url = server.URL + "/user/" + strconv.FormatInt(userID, 10) + "/cart"
	req, err = http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCartAPI_ClearCart(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.Close()

	// Сначала добавим товар в корзину
	userID := int64(1)
	sku := uint32(123)
	count := uint16(2)

	// Add item
	reqBody := map[string]uint16{
		"count": count,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	url := server.URL + "/user/" + strconv.FormatInt(userID, 10) + "/cart/" + strconv.FormatUint(uint64(sku), 10)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()

	// Clear cart
	url = server.URL + "/user/" + strconv.FormatInt(userID, 10) + "/cart"
	req, err = http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify cart is empty
	url = server.URL + "/user/" + strconv.FormatInt(userID, 10) + "/cart"
	req, err = http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
