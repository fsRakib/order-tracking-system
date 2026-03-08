package aggregate

import (
	"errors"

	"order-service/domain/event"
	"order-service/domain/valueobject"
)

// Order is the Aggregate Root for the order domain
// An Aggregate Root is the main entity that controls a cluster of related objects
// All state changes to an Order must go through these methods - no direct field access
//
// Think of it like a Java class where all fields are private and
// business rules are enforced in setters/methods
type Order struct {
	id           valueobject.OrderID
	customerID   valueobject.CustomerID
	customerName string
	items        []OrderItem
	status       valueobject.OrderStatus
	totalAmount  valueobject.Money

	// domainEvents holds events that happened to this order
	// They will be published to RabbitMQ by the application layer after saving
	domainEvents []event.DomainEvent
}

// NewOrder creates a new Order aggregate in Pending status
// This is the only way to create an Order - validates all required data upfront
//
// Java equivalent: public static Order create(String customerID) { ... }
func NewOrder(customerID valueobject.CustomerID, customerName string) (*Order, error) {
	if customerID.IsEmpty() {
		return nil, errors.New("customer ID is required to create an order")
	}

	// Start with zero money
	zeroMoney, err := valueobject.NewMoney(0, "USD")
	if err != nil {
		return nil, err
	}

	return &Order{
		id:           valueobject.NewOrderID(),
		customerID:   customerID,
		customerName: customerName,
		items:        []OrderItem{},
		status:       valueobject.StatusPending,
		totalAmount:  zeroMoney,
		domainEvents: []event.DomainEvent{},
	}, nil
}

// ReconstructOrder rebuilds an Order aggregate from database data
// Used by the repository when loading an existing order from Postgres
// Different from NewOrder because we don't generate a new ID or set Pending status
func ReconstructOrder(
	id valueobject.OrderID,
	customerID valueobject.CustomerID,
	customerName string,
	items []OrderItem,
	status valueobject.OrderStatus,
	totalAmount valueobject.Money,
) *Order {
	return &Order{
		id:           id,
		customerID:   customerID,
		customerName: customerName,
		items:        items,
		status:       status,
		totalAmount:  totalAmount,
		domainEvents: []event.DomainEvent{},
	}
}

// ---- Business Methods (where business rules live) ----

// AddItem adds a product to the order
// Business rules enforced:
//  1. Can only add items while order is still pending
//  2. Quantity must be positive
//  3. Maximum 10 items per order
func (o *Order) AddItem(sku valueobject.SKU, quantity int32, unitPrice valueobject.Money) error {
	if o.status != valueobject.StatusPending {
		return errors.New("items can only be added to a pending order")
	}

	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if len(o.items) >= 10 {
		return errors.New("an order cannot have more than 10 items")
	}

	// Check if same SKU already exists - if so, increase quantity instead
	for i, existing := range o.items {
		if existing.sku.Equals(sku) {
			o.items[i].quantity += quantity
			o.recalculateTotal()
			return nil
		}
	}

	// New item - append it
	newItem := newOrderItem(sku, quantity, unitPrice)
	o.items = append(o.items, newItem)
	o.recalculateTotal()

	return nil
}

// Confirm transitions order from Pending to Confirmed
// Business rule: only a Pending order can be confirmed
func (o *Order) Confirm() error {
	if !o.status.CanTransitionTo(valueobject.StatusConfirmed) {
		return errors.New("order cannot be confirmed from status: " + o.status.String())
	}

	if len(o.items) == 0 {
		return errors.New("cannot confirm an order with no items")
	}

	o.status = valueobject.StatusConfirmed

	// Record domain event - will be published to RabbitMQ later
	o.domainEvents = append(o.domainEvents, event.OrderCreatedEvent{
		OrderID:      o.id.Value(),
		CustomerID:   o.customerID.Value(),
		CustomerName: o.customerName,
		Status:       o.status.String(),
		TotalAmount:  o.totalAmount.Amount(),
		Items:        buildEventItems(o.items),
	})

	return nil
}

// UpdateStatus transitions the order to a new status
// Business rule: validates status transition is allowed
func (o *Order) UpdateStatus(newStatus valueobject.OrderStatus) error {
	if !o.status.CanTransitionTo(newStatus) {
		return errors.New("invalid status transition from " + o.status.String() + " to " + newStatus.String())
	}

	o.status = newStatus

	// Record domain event for status update
	o.domainEvents = append(o.domainEvents, event.OrderStatusUpdatedEvent{
		OrderID:      o.id.Value(),
		CustomerID:   o.customerID.Value(),
		CustomerName: o.customerName,
		Status:       o.status.String(),
		TotalAmount:  o.totalAmount.Amount(),
		Items:        buildEventItems(o.items),
	})

	return nil
}

// Cancel cancels the order
// Business rule: only pending or confirmed orders can be cancelled
func (o *Order) Cancel() error {
	if !o.status.CanTransitionTo(valueobject.StatusCancelled) {
		return errors.New("order cannot be cancelled from status: " + o.status.String())
	}

	o.status = valueobject.StatusCancelled
	return nil
}

// ---- Private Helpers ----

// recalculateTotal recalculates totalAmount from all items
// Called automatically after any item modification
func (o *Order) recalculateTotal() {
	total, _ := valueobject.NewMoney(0, "USD")
	for _, item := range o.items {
		subtotal := item.Subtotal()
		total, _ = total.Add(subtotal)
	}
	o.totalAmount = total
}

// buildEventItems converts OrderItems to event format
func buildEventItems(items []OrderItem) []event.OrderItemEvent {
	eventItems := make([]event.OrderItemEvent, len(items))
	for i, item := range items {
		eventItems[i] = event.OrderItemEvent{
			SKU:       item.sku.Value(),
			Quantity:  item.quantity,
			UnitPrice: item.unitPrice.Amount(),
		}
	}
	return eventItems
}

// ---- Getters (Public Read Access) ----

// ID returns the order's unique identifier
func (o *Order) ID() valueobject.OrderID { return o.id }

// CustomerID returns the customer who placed the order
func (o *Order) CustomerID() valueobject.CustomerID { return o.customerID }

// CustomerName returns the customer's display name
func (o *Order) CustomerName() string { return o.customerName }

// Items returns a copy of the order items (copy prevents external mutation)
func (o *Order) Items() []OrderItem {
	itemsCopy := make([]OrderItem, len(o.items))
	copy(itemsCopy, o.items)
	return itemsCopy
}

// Status returns the current order status
func (o *Order) Status() valueobject.OrderStatus { return o.status }

// TotalAmount returns the calculated total price
func (o *Order) TotalAmount() valueobject.Money { return o.totalAmount }

// DomainEvents returns all events raised during this operation
// The application layer reads these and publishes them to RabbitMQ
func (o *Order) DomainEvents() []event.DomainEvent { return o.domainEvents }

// ClearDomainEvents clears events after they have been published
func (o *Order) ClearDomainEvents() { o.domainEvents = []event.DomainEvent{} }
