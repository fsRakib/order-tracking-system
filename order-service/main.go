package main

import (
	"log"
	"net"
	"os"

	"order-service/application/command"
	"order-service/application/query"
	"order-service/db"
	grpcclient "order-service/infrastructure/grpc_client"
	"order-service/infrastructure/messaging"
	"order-service/infrastructure/persistence"
	grpchandler "order-service/interface/grpc"
	orderpb "order-service/pb/order"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	// ---- Infrastructure: Connect to external systems ----

	// Connect to PostgreSQL
	db.Connect()

	// Connect to Stock Service via gRPC
	stockServiceAddr := os.Getenv("STOCK_SERVICE_ADDR")
	stockService, err := grpcclient.NewStockGRPCClient(stockServiceAddr)
	if err != nil {
		log.Fatalf("failed to connect to stock service: %v", err)
	}

	// Connect to RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	rabbitConn, err := messaging.Connect(rabbitURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	// ---- Infrastructure: Create implementations ----

	orderRepo := persistence.NewPostgresOrderRepository(db.DB)
	customerService := persistence.NewPostgresCustomerService(db.DB)
	eventPublisher := messaging.NewRabbitMQPublisher(rabbitConn)

	// ---- Application: Create use case handlers ----

	createOrderHandler := command.NewCreateOrderHandler(
		orderRepo,
		stockService,
		customerService,
		eventPublisher,
	)
	updateOrderStatusHandler := command.NewUpdateOrderStatusHandler(
		orderRepo,
		eventPublisher,
	)
	getOrderHandler := query.NewGetOrderHandler(orderRepo)
	getOrdersByCustomerHandler := query.NewGetOrdersByCustomerHandler(orderRepo)

	// ---- Interface: Create gRPC handler ----

	orderHandler := grpchandler.NewOrderHandler(
		createOrderHandler,
		updateOrderStatusHandler,
		getOrderHandler,
		getOrdersByCustomerHandler,
	)

	// ---- Start gRPC server ----

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(grpcServer, orderHandler)
	reflection.Register(grpcServer)

	log.Printf("order service running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
