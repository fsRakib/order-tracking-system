package event

// DomainEvent is the base interface all domain events implement
// An event represents something significant that happened in the business
//
// Java equivalent: interface DomainEvent { String getEventType(); }
type DomainEvent interface {
	EventType() string
}

// OrderItemEvent is a helper struct used inside events
// Represents a single item snapshot at the time the event was raised
type OrderItemEvent struct {
	SKU       string  `json:"sku"`
	Quantity  int32   `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

// ---- Concrete Domain Events ----

// OrderCreatedEvent is raised when a new order is confirmed
// Published to RabbitMQ so analytics-service and other consumers can react
type OrderCreatedEvent struct {
	OrderID      string           `json:"order_id"`
	CustomerID   string           `json:"customer_id"`
	CustomerName string           `json:"customer_name"`
	Status       string           `json:"status"`
	TotalAmount  float64          `json:"total_amount"`
	Items        []OrderItemEvent `json:"items"`
}

// EventType returns the event name (used for routing/logging)
func (e OrderCreatedEvent) EventType() string {
	return "order.created"
}

// OrderStatusUpdatedEvent is raised when an order's status changes
// Published to RabbitMQ so other services can react (e.g., analytics tracking)
type OrderStatusUpdatedEvent struct {
	OrderID      string           `json:"order_id"`
	CustomerID   string           `json:"customer_id"`
	CustomerName string           `json:"customer_name"`
	Status       string           `json:"status"`
	TotalAmount  float64          `json:"total_amount"`
	Items        []OrderItemEvent `json:"items"`
}

// EventType returns the event name
func (e OrderStatusUpdatedEvent) EventType() string {
	return "order.status_updated"
}
