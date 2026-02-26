package main

import (
	"log"
	"net"
	"os"

	"order-service/client"
	"order-service/db"
	orderpb "order-service/pb/order"
	"order-service/publisher"
	"order-service/server"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	// Connect to PostgreSQL
	db.Connect()

	// Connect to Stock Service via gRPC
	client.ConnectStockService()

	// Connect to RabbitMQ
	publisher.Connect()

	// Start gRPC server
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(grpcServer, &server.OrderServer{})
	reflection.Register(grpcServer)

	log.Printf("order service running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
