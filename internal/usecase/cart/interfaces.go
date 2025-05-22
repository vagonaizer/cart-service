package cart

import (
	"route256/cart/internal/domain/ports"
)

//go:generate minimock -i CartRepository -o ./mocks/cart_repository_mock.go -g
//go:generate minimock -i ProductService -o ./mocks/product_service_mock.go -g

type CartRepository interface {
	ports.CartRepository
}

type ProductService interface {
	ports.ProductService
}
