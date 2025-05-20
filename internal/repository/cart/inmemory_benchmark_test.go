package cart

import (
	"route256/cart/internal/domain/models"
	"route256/cart/internal/infrastructure/repository/inmemory"
	"testing"
)

func BenchmarkInMemoryCartRepository_AddItem(b *testing.B) {
	repo := inmemory.NewCartRepository()
	cart := &models.Cart{
		UserID:     1,
		Items:      make(models.ItemList, 0),
		TotalPrice: 0,
	}
	err := repo.CreateCart(cart)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cart.Items = append(cart.Items, models.Item{
			SKU:      uint32(i % (1 << 32)),
			Quantity: 1,
			Price:    1000,
		})
		cart.TotalPrice += 1000
		err := repo.SaveCart(cart)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInMemoryCartRepository_GetCart(b *testing.B) {
	repo := inmemory.NewCartRepository()
	cart := &models.Cart{
		UserID: 1,
		Items: models.ItemList{
			{
				SKU:      123,
				Quantity: 2,
				Price:    1000,
			},
		},
		TotalPrice: 2000,
	}
	err := repo.SaveCart(cart)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetCart(1)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInMemoryCartRepository_RemoveItem(b *testing.B) {
	repo := inmemory.NewCartRepository()
	cart := &models.Cart{
		UserID: 1,
		Items: models.ItemList{
			{
				SKU:      123,
				Quantity: 2,
				Price:    1000,
			},
		},
		TotalPrice: 2000,
	}
	err := repo.SaveCart(cart)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cart.Items = make(models.ItemList, 0)
		cart.TotalPrice = 0
		err := repo.SaveCart(cart)
		if err != nil {
			b.Fatal(err)
		}
	}
}
