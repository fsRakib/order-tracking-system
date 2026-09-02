package valueobject

import "errors"

// OrderStatus represents the status of an order in the system
// Using a custom type instead of string prevents typos and provides type safety
type OrderStatus string

// Valid order status values (like enum in Java/C++)
const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

// NewOrderStatus creates a validated OrderStatus
// Returns error if the status is not valid
func NewOrderStatus(status string) (OrderStatus, error) {
	s := OrderStatus(status)
	if !s.IsValid() {
		return "", errors.New("invalid order status: " + status)
	}
	return s, nil
}

// IsValid checks if the status is one of the allowed values
func (s OrderStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusConfirmed, StatusShipped, StatusDelivered, StatusCancelled:
		return true
	}
	return false
}

// String returns the string representation
// This allows fmt.Printf("%s", status) to work
func (s OrderStatus) String() string {
	return string(s)
}

// CanTransitionTo checks if this status can transition to another status
// Business rule: certain status transitions are not allowed
func (s OrderStatus) CanTransitionTo(newStatus OrderStatus) bool {
	// Cancelled orders cannot transition to anything else
	if s == StatusCancelled {
		return false
	}

	// Delivered orders can only be cancelled (for returns)
	if s == StatusDelivered {
		return newStatus == StatusCancelled
	}

	// Pending can go to confirmed or cancelled
	if s == StatusPending {
		return newStatus == StatusConfirmed || newStatus == StatusCancelled
	}

	// Confirmed can go to shipped or cancelled
	if s == StatusConfirmed {
		return newStatus == StatusShipped || newStatus == StatusCancelled
	}

	// Shipped can go to delivered or cancelled
	if s == StatusShipped {
		return newStatus == StatusDelivered || newStatus == StatusCancelled
	}

	return false
}
