package valueobject

import "errors"

// SKU (Stock Keeping Unit) is a type-safe wrapper for product identifiers
// Ensures SKUs are properly validated before use
type SKU struct {
	value string
}

// NewSKU creates a validated SKU
func NewSKU(sku string) (SKU, error) {
	if sku == "" {
		return SKU{}, errors.New("SKU cannot be empty")
	}

	// Add validation rules for your SKU format
	// For example: minimum length, allowed characters, etc.
	if len(sku) < 3 {
		return SKU{}, errors.New("SKU must be at least 3 characters")
	}

	if len(sku) > 50 {
		return SKU{}, errors.New("SKU too long")
	}

	return SKU{value: sku}, nil
}

// Value returns the underlying string value
func (s SKU) Value() string {
	return s.value
}

// String returns the string representation
func (s SKU) String() string {
	return s.value
}

// Equals checks if two SKUs are equal
func (s SKU) Equals(other SKU) bool {
	return s.value == other.value
}

// IsEmpty checks if the SKU is empty
func (s SKU) IsEmpty() bool {
	return s.value == ""
}
