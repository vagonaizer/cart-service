package loms

import (
	"context"
	"log"

	loms "route256/cart/api/protos/gen/loms"
	"route256/cart/internal/domain/ports"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	lomsClient loms.LOMSClient
}

// NewClient creates a new LOMS client
func NewClient(address string) (ports.LOMSClient, error) {
	log.Printf("Connecting to LOMS service at %s", address)
	conn, err := grpc.DialContext(context.Background(), address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &client{
		lomsClient: loms.NewLOMSClient(conn),
	}, nil
}

func (c *client) CreateOrder(ctx context.Context, userID int64, items []ports.Item) (int64, error) {
	log.Printf("Creating order for user %d with %d items", userID, len(items))
	reqItems := make([]*loms.Item, len(items))
	for i, item := range items {
		reqItems[i] = &loms.Item{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		}
	}

	resp, err := c.lomsClient.OrderCreate(ctx, &loms.OrderCreateRequest{
		User:  userID,
		Items: reqItems,
	})
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return 0, err
	}

	log.Printf("Order created successfully with ID: %d", resp.OrderID)
	return resp.OrderID, nil
}

func (c *client) GetStocksInfo(ctx context.Context, sku uint32) (uint64, error) {
	log.Printf("Getting stock info for SKU %d", sku)
	resp, err := c.lomsClient.StocksInfo(ctx, &loms.StocksInfoRequest{
		Sku: sku,
	})
	if err != nil {
		log.Printf("Failed to get stock info: %v", err)
		return 0, err
	}

	log.Printf("Stock info for SKU %d: %d items available", sku, resp.Count)
	return resp.Count, nil
}

func (c *client) GetOrderInfo(ctx context.Context, orderID int64) (*ports.OrderInfo, error) {
	log.Printf("Getting order info for order %d", orderID)
	resp, err := c.lomsClient.OrderInfo(ctx, &loms.OrderInfoRequest{
		OrderID: orderID,
	})
	if err != nil {
		log.Printf("Failed to get order info: %v", err)
		return nil, err
	}

	items := make([]ports.Item, len(resp.Items))
	for i, item := range resp.Items {
		items[i] = ports.Item{
			SKU:   item.Sku,
			Count: uint16(item.Count),
		}
	}

	return &ports.OrderInfo{
		Status: resp.Status,
		UserID: resp.User,
		Items:  items,
	}, nil
}
