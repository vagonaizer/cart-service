.PHONY: build run

build:
	go build -o bin/cart-service ./cmd/cart

run: build
	./bin/cart-service
