package query

import (
	"context"
	"fmt"

	"order-service/application/dto"
	"order-service/domain/aggregate"
	"order-service/domain/repository"
	"order-service/domain/valueobject"
)

// GetOrderQuery carries the input needed to retrieve a single order
// Named "Query" because it only reads - does NOT change any state
type GetOrderQuery struct {
	OrderID string
}

// GetOrderHandler handles the GetOrder use case (read-only)
type GetOrderHandler struct {
	orderRepo repository.OrderRepository
}

// NewGetOrderHandler creates a new handler
func NewGetOrderHandler(orderRepo repository.OrderRepository) *GetOrderHandler {
	return &GetOrderHandler{orderRepo: orderRepo}
}

// Handle fetches a single order by ID and returns it as a DTO
func (h *GetOrderHandler) Handle(ctx context.Context, qry GetOrderQuery) (*dto.OrderDTO, error) {
	orderID, err := valueobject.NewOrderIDFromString(qry.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	order, err := h.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return toDTO(order), nil
}

// GetOrdersByCustomerQuery carries the input needed to list all orders of a customer
type GetOrdersByCustomerQuery struct {
	CustomerID string
}

// GetOrdersByCustomerHandler handles the GetOrdersByCustomer use case (read-only)
type GetOrdersByCustomerHandler struct {
	orderRepo repository.OrderRepository
}

// NewGetOrdersByCustomerHandler creates a new handler
func NewGetOrdersByCustomerHandler(orderRepo repository.OrderRepository) *GetOrdersByCustomerHandler {
	return &GetOrdersByCustomerHandler{orderRepo: orderRepo}
}

// Handle fetches all orders for a customer and returns them as DTOs
func (h *GetOrdersByCustomerHandler) Handle(ctx context.Context, qry GetOrdersByCustomerQuery) ([]*dto.OrderDTO, error) {
	customerID, err := valueobject.NewCustomerID(qry.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %w", err)
	}

	orders, err := h.orderRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.OrderDTO, 0, len(orders))
	for _, order := range orders {
		result = append(result, toDTO(order))
	}

	return result, nil
}

// toDTO converts an Order aggregate to an OrderDTO
// Defined here so both query handlers can use it
func toDTO(order *aggregate.Order) *dto.OrderDTO {
	items := make([]dto.OrderItemDTO, 0, len(order.Items()))
	for _, item := range order.Items() {
		items = append(items, dto.OrderItemDTO{
			SKU:       item.SKU().Value(),
			Quantity:  item.Quantity(),
			UnitPrice: item.UnitPrice().Amount(),
			Subtotal:  item.Subtotal().Amount(),
		})
	}

	return &dto.OrderDTO{
		OrderID:      order.ID().Value(),
		CustomerID:   order.CustomerID().Value(),
		CustomerName: order.CustomerName(),
		Status:       order.Status().String(),
		TotalAmount:  order.TotalAmount().Amount(),
		Items:        items,
	}
}
