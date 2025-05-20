package cart

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"route256/cart/internal/domain/models"
	"route256/cart/internal/usecase/cart/mocks"
)

func TestCartService_AddItem(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	repo := mocks.NewCartRepositoryMock(mc)
	productService := mocks.NewProductServiceMock(mc)
	service := NewCartService(repo, productService)

	t.Run("success", func(t *testing.T) {
		userID := int64(1)
		sku := uint32(123)
		count := uint16(2)

		product := &models.Product{
			SKU:   sku,
			Name:  "Test Product",
			Price: 1000,
		}

		productService.GetProductMock.Expect(sku).Return(product, nil)
		repo.GetCartMock.Expect(userID).Return(nil, models.ErrCartNotFound)
		repo.CreateCartMock.Set(func(cart *models.Cart) error {
			assert.Equal(t, userID, cart.UserID)
			return nil
		})
		repo.SaveCartMock.Set(func(cart *models.Cart) error {
			assert.Equal(t, userID, cart.UserID)
			assert.Len(t, cart.Items, 1)
			assert.Equal(t, sku, cart.Items[0].SKU)
			assert.Equal(t, count, cart.Items[0].Quantity)
			assert.Equal(t, product.Price, cart.Items[0].Price)
			assert.Equal(t, product.Price*uint32(count), cart.TotalPrice)
			return nil
		})

		err := service.AddItem(userID, sku, count)
		require.NoError(t, err)
	})

	t.Run("product not found", func(t *testing.T) {
		userID := int64(1)
		sku := uint32(123)
		count := uint16(2)

		productService.GetProductMock.Expect(sku).Return(nil, models.ErrProductNotFound)

		err := service.AddItem(userID, sku, count)
		assert.ErrorIs(t, err, models.ErrProductNotFound)
	})
}

func TestCartService_RemoveItem(t *testing.T) {
	tests := []struct {
		name    string
		userID  int64
		sku     uint32
		mock    func(mc *minimock.Controller) (CartRepository, ProductService)
		wantErr bool
	}{
		{
			name:   "successful remove item",
			userID: 1,
			sku:    123,
			mock: func(mc *minimock.Controller) (CartRepository, ProductService) {
				repo := mocks.NewCartRepositoryMock(mc)
				productService := mocks.NewProductServiceMock(mc)

				repo.GetCartMock.Expect(1).Return(&models.Cart{
					UserID: 1,
					Items: models.ItemList{
						{
							SKU:      123,
							Quantity: 2,
							Price:    1000,
						},
					},
					TotalPrice: 2000,
				}, nil)

				repo.SaveCartMock.Expect(&models.Cart{
					UserID:     1,
					Items:      make(models.ItemList, 0),
					TotalPrice: 0,
				}).Return(nil)

				return repo, productService
			},
			wantErr: false,
		},
		{
			name:   "cart not found",
			userID: 1,
			sku:    123,
			mock: func(mc *minimock.Controller) (CartRepository, ProductService) {
				repo := mocks.NewCartRepositoryMock(mc)
				productService := mocks.NewProductServiceMock(mc)

				repo.GetCartMock.Expect(1).Return(nil, models.ErrCartNotFound)

				return repo, productService
			},
			wantErr: false, // As per spec, return success if cart doesn't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			repo, productService := tt.mock(mc)
			service := NewCartService(repo, productService)

			err := service.RemoveItem(tt.userID, tt.sku)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCartService_ClearCart(t *testing.T) {
	tests := []struct {
		name    string
		userID  int64
		mock    func(mc *minimock.Controller) (CartRepository, ProductService)
		wantErr bool
	}{
		{
			name:   "successful clear cart",
			userID: 1,
			mock: func(mc *minimock.Controller) (CartRepository, ProductService) {
				repo := mocks.NewCartRepositoryMock(mc)
				productService := mocks.NewProductServiceMock(mc)

				repo.GetCartMock.Expect(1).Return(&models.Cart{
					UserID: 1,
					Items: models.ItemList{
						{
							SKU:      123,
							Quantity: 2,
							Price:    1000,
						},
					},
					TotalPrice: 2000,
				}, nil)

				repo.SaveCartMock.Expect(&models.Cart{
					UserID:     1,
					Items:      make(models.ItemList, 0),
					TotalPrice: 0,
				}).Return(nil)

				return repo, productService
			},
			wantErr: false,
		},
		{
			name:   "cart not found",
			userID: 1,
			mock: func(mc *minimock.Controller) (CartRepository, ProductService) {
				repo := mocks.NewCartRepositoryMock(mc)
				productService := mocks.NewProductServiceMock(mc)

				repo.GetCartMock.Expect(1).Return(nil, models.ErrCartNotFound)

				return repo, productService
			},
			wantErr: false, // As per spec, return success if cart doesn't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			repo, productService := tt.mock(mc)
			service := NewCartService(repo, productService)

			err := service.ClearCart(tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCartService_GetCart(t *testing.T) {
	tests := []struct {
		name    string
		userID  int64
		mock    func(mc *minimock.Controller) (CartRepository, ProductService)
		want    *models.Cart
		wantErr bool
		errType error
	}{
		{
			name:   "successful get cart",
			userID: 1,
			mock: func(mc *minimock.Controller) (CartRepository, ProductService) {
				repo := mocks.NewCartRepositoryMock(mc)
				productService := mocks.NewProductServiceMock(mc)

				repo.GetCartMock.Expect(1).Return(&models.Cart{
					UserID: 1,
					Items: models.ItemList{
						{
							SKU:      123,
							Quantity: 2,
							Price:    1000,
						},
					},
					TotalPrice: 2000,
				}, nil)

				return repo, productService
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
			wantErr: false,
		},
		{
			name:   "cart not found",
			userID: 1,
			mock: func(mc *minimock.Controller) (CartRepository, ProductService) {
				repo := mocks.NewCartRepositoryMock(mc)
				productService := mocks.NewProductServiceMock(mc)

				repo.GetCartMock.Expect(1).Return(nil, models.ErrCartNotFound)

				return repo, productService
			},
			wantErr: true,
			errType: models.ErrCartNotFound,
		},
		{
			name:   "empty cart",
			userID: 1,
			mock: func(mc *minimock.Controller) (CartRepository, ProductService) {
				repo := mocks.NewCartRepositoryMock(mc)
				productService := mocks.NewProductServiceMock(mc)

				repo.GetCartMock.Expect(1).Return(&models.Cart{
					UserID:     1,
					Items:      make(models.ItemList, 0),
					TotalPrice: 0,
				}, nil)

				return repo, productService
			},
			wantErr: true,
			errType: models.ErrCartNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			repo, productService := tt.mock(mc)
			service := NewCartService(repo, productService)

			got, err := service.GetCart(tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.errType)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
