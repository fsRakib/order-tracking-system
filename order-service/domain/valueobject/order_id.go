package valueobject

import (
	"errors"

	"github.com/google/uuid"
)

// OrderID is a type-safe wrapper for order identifiers
// Prevents mixing up different types of IDs
type OrderID struct {
	value string
}

// NewOrderID creates a new random OrderID using UUID
func NewOrderID() OrderID {
	return OrderID{value: uuid.New().String()}
}

// NewOrderIDFromString creates an OrderID from an existing string
// Used when loading orders from database
func NewOrderIDFromString(id string) (OrderID, error) {
	if id == "" {
		return OrderID{}, errors.New("order ID cannot be empty")
	}

	// Validate that it's a valid UUID format
	if _, err := uuid.Parse(id); err != nil {
		return OrderID{}, errors.New("invalid order ID format")
	}

	return OrderID{value: id}, nil
}

// Value returns the underlying string value
func (o OrderID) Value() string {
	return o.value
}

// String returns the string representation
func (o OrderID) String() string {
	return o.value
}

// Equals checks if two OrderIDs are equal
func (o OrderID) Equals(other OrderID) bool {
	return o.value == other.value
}

// IsEmpty checks if the OrderID is empty
func (o OrderID) IsEmpty() bool {
	return o.value == ""
}
