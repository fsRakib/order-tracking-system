package server

import (
	"context"
	"log"

	"stock-service/db"
	pb "stock-service/pb/stock"
)

// StockServer implements the generated StockServiceServer interface
type StockServer struct {
	pb.UnimplementedStockServiceServer
}

func (s *StockServer) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {
	log.Printf("GetStock called for SKU: %s", req.Sku)

	stock, err := db.GetStockBySKU(req.Sku)
	if err != nil {
		return nil, err
	}

	return &pb.GetStockResponse{
		Sku:      stock.SKU,
		Quantity: int32(stock.Quantity),
	}, nil
}

func (s *StockServer) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	log.Printf("ReserveStock called: SKU=%s, Quantity=%d", req.Sku, req.Quantity)

	err := db.ReserveStock(req.Sku, req.Quantity)
	if err != nil {
		log.Printf("ReserveStock failed: %v", err)
		return &pb.ReserveStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ReserveStockResponse{
		Success: true,
		Message: "stock reserved successfully",
	}, nil
}

func (s *StockServer) ReleaseStock(ctx context.Context, req *pb.ReleaseStockRequest) (*pb.ReleaseStockResponse, error) {
	log.Printf("ReleaseStock called: SKU=%s, Quantity=%d", req.Sku, req.Quantity)

	err := db.ReleaseStock(req.Sku, req.Quantity)
	if err != nil {
		log.Printf("ReleaseStock failed: %v", err)
		return &pb.ReleaseStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ReleaseStockResponse{
		Success: true,
		Message: "stock released successfully",
	}, nil
}