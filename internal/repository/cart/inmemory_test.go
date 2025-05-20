package cart

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"route256/cart/internal/domain/models"
	"route256/cart/internal/infrastructure/repository/inmemory"
)

func TestInMemoryCartRepository_GetCart(t *testing.T) {
	tests := []struct {
		name    string
		userID  int64
		setup   func(repo *inmemory.CartRepository)
		want    *models.Cart
		wantErr error
	}{
		{
			name:   "cart exists",
			userID: 1,
			setup: func(repo *inmemory.CartRepository) {
				_ = repo.SaveCart(&models.Cart{
					UserID: 1,
					Items: models.ItemList{
						{
							SKU:      123,
							Quantity: 2,
							Price:    1000,
						},
					},
					TotalPrice: 2000,
				})
			},
			want: &models.Cart{
				UserID: 1,
				Items: models.ItemList{
					{
						SKU:      123,
						Quantity: 2,
						Price:    1000,
					},
				},
				TotalPrice: 2000,
			},
			wantErr: nil,
		},
		{
			name:    "cart not found",
			userID:  1,
			setup:   func(repo *inmemory.CartRepository) {},
			want:    nil,
			wantErr: models.ErrCartNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inmemory.NewCartRepository()
			tt.setup(repo)

			got, err := repo.GetCart(tt.userID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInMemoryCartRepository_SaveCart(t *testing.T) {
	tests := []struct {
		name    string
		cart    *models.Cart
		setup   func(repo *inmemory.CartRepository)
		wantErr error
	}{
		{
			name: "save new cart",
			cart: &models.Cart{
				UserID: 1,
				Items: models.ItemList{
					{
						SKU:      123,
						Quantity: 2,
						Price:    1000,
					},
				},
				TotalPrice: 2000,
			},
			setup:   func(repo *inmemory.CartRepository) {},
			wantErr: nil,
		},
		{
			name: "update existing cart",
			cart: &models.Cart{
				UserID: 1,
				Items: models.ItemList{
					{
						SKU:      123,
						Quantity: 3,
						Price:    1000,
					},
				},
				TotalPrice: 3000,
			},
			setup: func(repo *inmemory.CartRepository) {
				_ = repo.SaveCart(&models.Cart{
					UserID: 1,
					Items: models.ItemList{
						{
							SKU:      123,
							Quantity: 2,
							Price:    1000,
						},
					},
					TotalPrice: 2000,
				})
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inmemory.NewCartRepository()
			tt.setup(repo)

			err := repo.SaveCart(tt.cart)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				got, err := repo.GetCart(tt.cart.UserID)
				require.NoError(t, err)
				assert.Equal(t, tt.cart, got)
			}
		})
	}
}

func TestInMemoryCartRepository_CreateCart(t *testing.T) {
	tests := []struct {
		name    string
		cart    *models.Cart
		setup   func(repo *inmemory.CartRepository)
		wantErr error
	}{
		{
			name: "create new cart",
			cart: &models.Cart{
				UserID:     1,
				Items:      make(models.ItemList, 0),
				TotalPrice: 0,
			},
			setup:   func(repo *inmemory.CartRepository) {},
			wantErr: nil,
		},
		{
			name: "cart already exists",
			cart: &models.Cart{
				UserID:     1,
				Items:      make(models.ItemList, 0),
				TotalPrice: 0,
			},
			setup: func(repo *inmemory.CartRepository) {
				_ = repo.SaveCart(&models.Cart{
					UserID:     1,
					Items:      make(models.ItemList, 0),
					TotalPrice: 0,
				})
			},
			wantErr: models.ErrCartAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inmemory.NewCartRepository()
			tt.setup(repo)

			err := repo.CreateCart(tt.cart)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				got, err := repo.GetCart(tt.cart.UserID)
				require.NoError(t, err)
				assert.Equal(t, tt.cart, got)
			}
		})
	}
}
