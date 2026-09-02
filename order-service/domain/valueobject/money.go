package valueobject

import (
	"errors"
	"fmt"
)

// Money represents monetary value with precision
// Stores amount in cents to avoid float precision issues
// Similar to Java's BigDecimal or C++ decimal libraries
type Money struct {
	amountInCents int64  // Store as cents (e.g., $10.50 = 1050 cents)
	currency      string // Currency code like "USD", "EUR"
}

// NewMoney creates a Money value object from a float amount
// Example: NewMoney(10.50, "USD") creates $10.50
func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}

	if currency == "" {
		currency = "USD" // Default currency
	}

	// Convert dollars to cents (multiply by 100)
	// This avoids floating-point precision issues
	amountInCents := int64(amount * 100)

	return Money{
		amountInCents: amountInCents,
		currency:      currency,
	}, nil
}

// NewMoneyFromCents creates Money directly from cents
// Useful when you already have the cent value
func NewMoneyFromCents(cents int64, currency string) (Money, error) {
	if cents < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}

	if currency == "" {
		currency = "USD"
	}

	return Money{
		amountInCents: cents,
		currency:      currency,
	}, nil
}

// Amount returns the money value as a float (in dollars)
// Example: If amountInCents = 1050, returns 10.50
func (m Money) Amount() float64 {
	return float64(m.amountInCents) / 100.0
}

// AmountInCents returns the raw cent value
func (m Money) AmountInCents() int64 {
	return m.amountInCents
}

// Currency returns the currency code
func (m Money) Currency() string {
	return m.currency
}

// Add adds two Money values together
// Returns error if currencies don't match
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot add %s and %s", m.currency, other.currency)
	}

	return Money{
		amountInCents: m.amountInCents + other.amountInCents,
		currency:      m.currency,
	}, nil
}

// Subtract subtracts other money from this money
// Returns error if currencies don't match or result would be negative
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot subtract %s from %s", other.currency, m.currency)
	}

	if m.amountInCents < other.amountInCents {
		return Money{}, errors.New("result would be negative")
	}

	return Money{
		amountInCents: m.amountInCents - other.amountInCents,
		currency:      m.currency,
	}, nil
}

// Multiply multiplies the money by a factor
// Example: Money(10.00).Multiply(3) = 30.00
func (m Money) Multiply(factor int32) Money {
	return Money{
		amountInCents: m.amountInCents * int64(factor),
		currency:      m.currency,
	}
}

// IsZero checks if the amount is zero
func (m Money) IsZero() bool {
	return m.amountInCents == 0
}

// IsPositive checks if the amount is greater than zero
func (m Money) IsPositive() bool {
	return m.amountInCents > 0
}

// Equals checks if two Money values are equal
func (m Money) Equals(other Money) bool {
	return m.amountInCents == other.amountInCents && m.currency == other.currency
}

// String returns a formatted string representation
// Example: "$10.50 USD"
func (m Money) String() string {
	return fmt.Sprintf("$%.2f %s", m.Amount(), m.currency)
}
