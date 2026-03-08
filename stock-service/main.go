package main

import (
	"log"
	"net"
	"os"

	"stock-service/consumer"
	"stock-service/db"
	pb "order-tracking-system/pb/stock"
	"stock-service/server"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	// Connect to database
	db.Connect()

	// Start RabbitMQ consumer in background
	go consumer.StartConsumer()

	// Start gRPC server
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStockServiceServer(grpcServer, &server.StockServer{})
	reflection.Register(grpcServer)

	log.Printf("stock service running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
