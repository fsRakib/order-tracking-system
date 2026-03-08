package command

import (
	"context"
	"fmt"
	"log"

	"order-service/application/dto"
	"order-service/application/port"
	"order-service/domain/repository"
	"order-service/domain/valueobject"
)

// UpdateOrderStatusCommand carries the input data needed to update an order's status
type UpdateOrderStatusCommand struct {
	OrderID   string
	NewStatus string
}

// UpdateOrderStatusHandler handles the UpdateOrderStatus use case
type UpdateOrderStatusHandler struct {
	orderRepo      repository.OrderRepository
	eventPublisher port.EventPublisher
}

// NewUpdateOrderStatusHandler creates a new handler with required dependencies
func NewUpdateOrderStatusHandler(
	orderRepo repository.OrderRepository,
	eventPublisher port.EventPublisher,
) *UpdateOrderStatusHandler {
	return &UpdateOrderStatusHandler{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle executes the UpdateOrderStatus use case
// Returns the updated OrderDTO or an error
func (h *UpdateOrderStatusHandler) Handle(ctx context.Context, cmd UpdateOrderStatusCommand) (*dto.OrderDTO, error) {
	// Step 1: Validate and create typed value objects
	orderID, err := valueobject.NewOrderIDFromString(cmd.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	newStatus, err := valueobject.NewOrderStatus(cmd.NewStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid order status: %w", err)
	}

	// Step 2: Load the existing order from the repository
	order, err := h.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Step 3: Apply status change through the aggregate
	// The aggregate enforces valid transitions (e.g., cannot go from Shipped back to Pending)
	if err := order.UpdateStatus(newStatus); err != nil {
		return nil, fmt.Errorf("invalid status transition: %w", err)
	}

	// Step 4: Persist the updated order
	if err := h.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	log.Printf("order status updated: %s → %s", cmd.OrderID, cmd.NewStatus)

	// Step 5: Publish domain events raised by the aggregate (fire-and-forget)
	go func() {
		for _, domainEvent := range order.DomainEvents() {
			if err := h.eventPublisher.Publish(context.Background(), domainEvent); err != nil {
				log.Printf("warning: failed to publish event %s: %v", domainEvent.EventType(), err)
			}
		}
		order.ClearDomainEvents()
	}()

	return toDTO(order), nil
}
