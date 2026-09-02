package command

import (
	"context"
	"fmt"
	"log"

	"order-service/application/dto"
	"order-service/application/port"
	"order-service/domain/aggregate"
	"order-service/domain/repository"
	"order-service/domain/valueobject"
)

// CreateOrderCommand carries the input data needed to create an order
// Named "Command" because it changes system state (write operation)
// In Java this would be: class CreateOrderCommand { String customerId; List<Item> items; }
type CreateOrderCommand struct {
	CustomerID string
	Items      []CreateOrderItemCommand
}

// CreateOrderItemCommand carries data for a single item in the command
type CreateOrderItemCommand struct {
	SKU       string
	Quantity  int32
	UnitPrice float64
}

// CreateOrderHandler handles the CreateOrder use case
// It orchestrates: stock reservation → aggregate creation → persistence → event publishing
//
// Dependencies are injected via the constructor (like Spring @Autowired in Java)
// All dependencies are interfaces - no concrete types
type CreateOrderHandler struct {
	orderRepo      repository.OrderRepository
	stockService   port.StockService
	customerService port.CustomerService
	eventPublisher port.EventPublisher
}

// NewCreateOrderHandler creates a new handler with all required dependencies
func NewCreateOrderHandler(
	orderRepo repository.OrderRepository,
	stockService port.StockService,
	customerService port.CustomerService,
	eventPublisher port.EventPublisher,
) *CreateOrderHandler {
	return &CreateOrderHandler{
		orderRepo:      orderRepo,
		stockService:   stockService,
		customerService: customerService,
		eventPublisher: eventPublisher,
	}
}

// Handle executes the CreateOrder use case step by step
// Returns an OrderDTO with the created order details, or an error
func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (*dto.OrderDTO, error) {
	// Step 1: Validate input
	if len(cmd.Items) == 0 {
		return nil, fmt.Errorf("order must have at least one item")
	}

	// Step 2: Reserve stock for each item via stock service
	// If any item fails, the loop returns early - stock already reserved is NOT rolled back here
	// (In production you would use a saga pattern for this, but that is advanced)
	for _, item := range cmd.Items {
		if err := h.stockService.ReserveStock(ctx, item.SKU, item.Quantity); err != nil {
			return nil, fmt.Errorf("stock reservation failed: %w", err)
		}
	}

	// Step 3: Fetch customer name for event publishing
	customerName, err := h.customerService.GetCustomerName(ctx, cmd.CustomerID)
	if err != nil {
		// Non-fatal: continue with empty name if customer lookup fails
		log.Printf("warning: could not fetch customer name for %s: %v", cmd.CustomerID, err)
		customerName = ""
	}

	// Step 4: Create customer ID value object (validates it is non-empty)
	customerID, err := valueobject.NewCustomerID(cmd.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %w", err)
	}

	// Step 5: Create Order aggregate - business rules enforced inside aggregate
	order, err := aggregate.NewOrder(customerID, customerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Step 6: Add items to the order aggregate
	// Each AddItem call enforces business rules (max 10. positive quantity, etc.)
	for _, item := range cmd.Items {
		sku, err := valueobject.NewSKU(item.SKU)
		if err != nil {
			return nil, fmt.Errorf("invalid SKU %s: %w", item.SKU, err)
		}

		unitPrice, err := valueobject.NewMoney(item.UnitPrice, "USD")
		if err != nil {
			return nil, fmt.Errorf("invalid unit price for SKU %s: %w", item.SKU, err)
		}

		if err := order.AddItem(sku, item.Quantity, unitPrice); err != nil {
			return nil, fmt.Errorf("failed to add item %s: %w", item.SKU, err)
		}
	}

	// Step 7: Confirm the order - transitions status Pending → Confirmed
	// Also raises the OrderCreatedEvent domain event
	if err := order.Confirm(); err != nil {
		return nil, fmt.Errorf("failed to confirm order: %w", err)
	}

	// Step 8: Persist the order via repository
	if err := h.orderRepo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	log.Printf("order saved: %s", order.ID().Value())

	// Step 9: Publish domain events collected by the aggregate (fire-and-forget)
	go func() {
		for _, domainEvent := range order.DomainEvents() {
			if err := h.eventPublisher.Publish(context.Background(), domainEvent); err != nil {
				log.Printf("warning: failed to publish event %s: %v", domainEvent.EventType(), err)
			}
		}
		order.ClearDomainEvents()
	}()

	// Step 10: Return DTO to the caller (gRPC handler)
	return toDTO(order), nil
}

// toDTO converts an Order aggregate to an OrderDTO for the caller
// This prevents the aggregate from leaking outside the application layer
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
