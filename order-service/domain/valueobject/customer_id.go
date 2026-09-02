package valueobject

import "errors"

// CustomerID is a type-safe wrapper for customer identifiers
// Using a custom type instead of string prevents accidentally mixing IDs
// (e.g., can't accidentally use OrderID where CustomerID is expected)
type CustomerID struct {
	value string
}

// NewCustomerID creates a validated CustomerID
func NewCustomerID(id string) (CustomerID, error) {
	if id == "" {
		return CustomerID{}, errors.New("customer ID cannot be empty")
	}

	// Add any other validation rules here
	// For example: length checks, format validation, etc.
	if len(id) > 100 {
		return CustomerID{}, errors.New("customer ID too long")
	}

	return CustomerID{value: id}, nil
}

// Value returns the underlying string value
func (c CustomerID) Value() string {
	return c.value
}

// String returns the string representation
func (c CustomerID) String() string {
	return c.value
}

// Equals checks if two CustomerIDs are equal
func (c CustomerID) Equals(other CustomerID) bool {
	return c.value == other.value
}

// IsEmpty checks if the CustomerID is empty
func (c CustomerID) IsEmpty() bool {
	return c.value == ""
}
