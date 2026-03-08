package grpc

import (
	"context"
	"log"

	"order-service/application/command"
	"order-service/application/dto"
	"order-service/application/query"
	orderpb "order-service/pb/order"
)

// OrderHandler is the gRPC server handler for the Order service
// It is intentionally thin - its ONLY job is:
//  1. Receive gRPC request
//  2. Convert to a Command or Query
//  3. Delegate to the application handler
//  4. Convert result DTO back to gRPC response
//
// No business logic, no SQL, no RabbitMQ, no gRPC client code here.
type OrderHandler struct {
	orderpb.UnimplementedOrderServiceServer

	// Command handlers (write operations)
	createOrderHandler      *command.CreateOrderHandler
	updateOrderStatusHandler *command.UpdateOrderStatusHandler

	// Query handlers (read operations)
	getOrderHandler            *query.GetOrderHandler
	getOrdersByCustomerHandler *query.GetOrdersByCustomerHandler
}

// NewOrderHandler creates a new gRPC handler with all application handlers injected
func NewOrderHandler(
	createOrder *command.CreateOrderHandler,
	updateStatus *command.UpdateOrderStatusHandler,
	getOrder *query.GetOrderHandler,
	getOrdersByCustomer *query.GetOrdersByCustomerHandler,
) *OrderHandler {
	return &OrderHandler{
		createOrderHandler:         createOrder,
		updateOrderStatusHandler:   updateStatus,
		getOrderHandler:            getOrder,
		getOrdersByCustomerHandler: getOrdersByCustomer,
	}
}

// CreateOrder receives a gRPC CreateOrder request and delegates to CreateOrderHandler
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	log.Printf("gRPC CreateOrder called for customer: %s", req.CustomerId)

	// Convert gRPC request → command
	items := make([]command.CreateOrderItemCommand, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, command.CreateOrderItemCommand{
			SKU:       item.Sku,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	cmd := command.CreateOrderCommand{
		CustomerID: req.CustomerId,
		Items:      items,
	}

	// Delegate to application layer
	result, err := h.createOrderHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Convert DTO → gRPC response
	return &orderpb.CreateOrderResponse{
		OrderId: result.OrderID,
		Status:  result.Status,
		Message: "order created successfully",
	}, nil
}

// GetOrder receives a gRPC GetOrder request and delegates to GetOrderHandler
func (h *OrderHandler) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.OrderResponse, error) {
	log.Printf("gRPC GetOrder called for order: %s", req.OrderId)

	result, err := h.getOrderHandler.Handle(ctx, query.GetOrderQuery{
		OrderID: req.OrderId,
	})
	if err != nil {
		return nil, err
	}

	return toProtoResponse(result), nil
}

// UpdateOrderStatus receives a gRPC UpdateOrderStatus request and delegates to UpdateOrderStatusHandler
func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.OrderResponse, error) {
	log.Printf("gRPC UpdateOrderStatus called: orderID=%s, status=%s", req.OrderId, req.Status)

	result, err := h.updateOrderStatusHandler.Handle(ctx, command.UpdateOrderStatusCommand{
		OrderID:   req.OrderId,
		NewStatus: req.Status,
	})
	if err != nil {
		return nil, err
	}

	return toProtoResponse(result), nil
}

// GetOrdersByCustomer receives a gRPC GetOrdersByCustomer request and delegates to the query handler
func (h *OrderHandler) GetOrdersByCustomer(ctx context.Context, req *orderpb.GetOrdersByCustomerRequest) (*orderpb.OrderListResponse, error) {
	log.Printf("gRPC GetOrdersByCustomer called for customer: %s", req.CustomerId)

	results, err := h.getOrdersByCustomerHandler.Handle(ctx, query.GetOrdersByCustomerQuery{
		CustomerID: req.CustomerId,
	})
	if err != nil {
		return nil, err
	}

	responses := make([]*orderpb.OrderResponse, 0, len(results))
	for _, r := range results {
		responses = append(responses, toProtoResponse(r))
	}

	return &orderpb.OrderListResponse{Orders: responses}, nil
}

// toProtoResponse converts an OrderDTO to a protobuf OrderResponse
func toProtoResponse(o *dto.OrderDTO) *orderpb.OrderResponse {
	items := make([]*orderpb.OrderItemResponse, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, &orderpb.OrderItemResponse{
			Sku:       item.SKU,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	return &orderpb.OrderResponse{
		OrderId:     o.OrderID,
		CustomerId:  o.CustomerID,
		Status:      o.Status,
		TotalAmount: o.TotalAmount,
		Items:       items,
	}
}
