.PHONY: build run test test-coverage lint bench install-tools generate-mocks

build:
	go build -o bin/cart-service ./cmd/cart

run: build
	./bin/cart-service

test:
	go test -v ./...

test-integration:
	go test -v ./... -tags=integration

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

lint:
	golangci-lint run

bench:
	go test -bench=. -benchmem ./...

install-tools:
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/uudashr/gocognit/cmd/gocognit@latest

generate-mocks:
	minimock -i ./internal/domain/ports/cart_service.go -o ./internal/usecase/cart/mocks/cart_service_mock.go
	minimock -i ./internal/domain/ports/product_service.go -o ./internal/usecase/cart/mocks/product_service_mock.go
	minimock -i ./internal/domain/ports/cart_repository.go -o ./internal/usecase/cart/mocks/cart_repository_mock.go
