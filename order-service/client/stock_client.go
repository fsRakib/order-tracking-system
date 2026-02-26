package client

import (
	"log"
	"os"

	pb "order-service/pb/stock"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var StockClient pb.StockServiceClient

func ConnectStockService() {
	addr := os.Getenv("STOCK_SERVICE_ADDR")
	if addr == "" {
		addr = "localhost:50052"
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to stock service: %v", err)
	}

	StockClient = pb.NewStockServiceClient(conn)
	log.Printf("connected to stock service at %s", addr)
}