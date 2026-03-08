package port

import "context"

// StockService defines how to interact with the Stock Service
//
// The application layer needs to reserve/release stock when processing orders
// but should NOT know about gRPC, proto files, or HTTP.
//
// The actual gRPC implementation is in infrastructure/grpc_client/
//
// Java equivalent:
//
//	public interface StockService {
//	    void reserveStock(String sku, int quantity) throws StockException;
//	    void releaseStock(String sku, int quantity) throws StockException;
//	}
type StockService interface {
	// ReserveStock reserves a quantity of a product in the stock service
	// Returns error if there is insufficient stock
	ReserveStock(ctx context.Context, sku string, quantity int32) error

	// ReleaseStock releases previously reserved stock
	// Called when an order is cancelled
	ReleaseStock(ctx context.Context, sku string, quantity int32) error
}
