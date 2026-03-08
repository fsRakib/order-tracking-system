package grpcclient

import (
	"context"
	"fmt"
	"log"

	"order-service/application/port"
	stockpb "order-tracking-system/pb/stock"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// StockGRPCClient implements port.StockService using gRPC
// This is the ONLY place in the codebase that knows about gRPC or proto files
//
// Java equivalent: class StockGRPCClient implements StockService { ... }
type StockGRPCClient struct {
	client stockpb.StockServiceClient
}

// NewStockGRPCClient creates a new gRPC client connected to the stock service
// Returns the interface type to enforce the contract
func NewStockGRPCClient(addr string) (port.StockService, error) {
	if addr == "" {
		addr = "localhost:50052"
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to stock service at %s: %w", addr, err)
	}

	log.Printf("connected to stock service at %s", addr)
	return &StockGRPCClient{
		client: stockpb.NewStockServiceClient(conn),
	}, nil
}

// ReserveStock calls the stock service to reserve inventory for an order item
// Returns error if stock is insufficient
func (c *StockGRPCClient) ReserveStock(ctx context.Context, sku string, quantity int32) error {
	resp, err := c.client.ReserveStock(ctx, &stockpb.ReserveStockRequest{
		Sku:      sku,
		Quantity: quantity,
	})
	if err != nil {
		return fmt.Errorf("stock service error for SKU %s: %w", sku, err)
	}

	if !resp.Success {
		return fmt.Errorf("insufficient stock for SKU %s: %s", sku, resp.Message)
	}

	log.Printf("stock reserved: SKU=%s quantity=%d", sku, quantity)
	return nil
}

// ReleaseStock calls the stock service to return reserved inventory
// Called when an order is cancelled
func (c *StockGRPCClient) ReleaseStock(ctx context.Context, sku string, quantity int32) error {
	resp, err := c.client.ReleaseStock(ctx, &stockpb.ReleaseStockRequest{
		Sku:      sku,
		Quantity: quantity,
	})
	if err != nil {
		return fmt.Errorf("stock service error releasing SKU %s: %w", sku, err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to release stock for SKU %s: %s", sku, resp.Message)
	}

	log.Printf("stock released: SKU=%s quantity=%d", sku, quantity)
	return nil
}
