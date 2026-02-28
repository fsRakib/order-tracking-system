package server

import (
	"context"
	"fmt"
	"log"

	"order-service/client"
	"order-service/db"
	orderpb "order-service/pb/order"
	stockpb "order-service/pb/stock"
	"order-service/publisher"

	"github.com/google/uuid"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	log.Printf("CreateOrder called for customer: %s", req.CustomerId)

	// Step 1: Reserve stock for each item via Stock Service gRPC
	for _, item := range req.Items {
		resp, err := client.StockClient.ReserveStock(ctx, &stockpb.ReserveStockRequest{
			Sku:      item.Sku,
			Quantity: item.Quantity,
		})
		if err != nil {
			return nil, fmt.Errorf("stock service error for SKU %s: %v", item.Sku, err)
		}
		if !resp.Success {
			return nil, fmt.Errorf("failed to reserve stock for SKU %s: %s", item.Sku, resp.Message)
		}
		log.Printf("stock reserved for SKU: %s, quantity: %d", item.Sku, item.Quantity)
	}

	// Step 2: Calculate total amount
	var totalAmount float64
	var dbItems []db.OrderItem
	for _, item := range req.Items {
		totalAmount += float64(item.Quantity) * item.UnitPrice
		dbItems = append(dbItems, db.OrderItem{
			SKU:       item.Sku,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	// Step 3: Save order to database
	orderID := uuid.New().String()
	order := db.Order{
		ID:          orderID,
		CustomerID:  req.CustomerId,
		TotalAmount: totalAmount,
		Status:      "confirmed",
		Items:       dbItems,
	}

	if err := db.InsertOrder(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %v", err)
	}
	log.Printf("order saved to database: %s", orderID)

	// Step 4: Fetch customer name
	customerName, err := db.GetCustomerName(req.CustomerId)
	if err != nil {
		log.Printf("warning: failed to fetch customer name: %v", err)
		customerName = "" // Continue with empty name if fetch fails
	}

	// Step 5: Publish OrderCreated event to RabbitMQ asynchronously
	go func() {
		var eventItems []publisher.OrderItem
		for _, item := range req.Items {
			eventItems = append(eventItems, publisher.OrderItem{
				SKU:       item.Sku,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			})
		}

		publisher.PublishOrderCreated(publisher.OrderEvent{
			OrderID:      orderID,
			CustomerID:   req.CustomerId,
			CustomerName: customerName,
			Status:       "confirmed",
			TotalAmount:  totalAmount,
			Items:        eventItems,
		})
	}()

	return &orderpb.CreateOrderResponse{
		OrderId: orderID,
		Status:  "confirmed",
		Message: "order created successfully",
	}, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.OrderResponse, error) {
	log.Printf("GetOrder called for order: %s", req.OrderId)

	order, err := db.GetOrderByID(req.OrderId)
	if err != nil {
		return nil, err
	}

	return toOrderResponse(order), nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.OrderResponse, error) {
	log.Printf("UpdateOrderStatus called: orderID=%s, status=%s", req.OrderId, req.Status)

	order, err := db.UpdateOrderStatus(req.OrderId, req.Status)
	if err != nil {
		return nil, err
	}

	return toOrderResponse(order), nil
}

func (s *OrderServer) GetOrdersByCustomer(ctx context.Context, req *orderpb.GetOrdersByCustomerRequest) (*orderpb.OrderListResponse, error) {
	log.Printf("GetOrdersByCustomer called for customer: %s", req.CustomerId)

	orders, err := db.GetOrdersByCustomerID(req.CustomerId)
	if err != nil {
		return nil, err
	}

	var responses []*orderpb.OrderResponse
	for _, o := range orders {
		order := o
		responses = append(responses, toOrderResponse(&order))
	}

	return &orderpb.OrderListResponse{Orders: responses}, nil
}

// toOrderResponse converts a db.Order to a protobuf OrderResponse
func toOrderResponse(o *db.Order) *orderpb.OrderResponse {
	var items []*orderpb.OrderItemResponse
	for _, item := range o.Items {
		items = append(items, &orderpb.OrderItemResponse{
			Sku:       item.SKU,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	return &orderpb.OrderResponse{
		OrderId:     o.ID,
		CustomerId:  o.CustomerID,
		Status:      o.Status,
		TotalAmount: o.TotalAmount,
		CreatedAt:   o.CreatedAt,
		Items:       items,
	}
}
