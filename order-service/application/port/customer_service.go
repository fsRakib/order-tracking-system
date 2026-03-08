package port

import "context"

// CustomerService defines how to retrieve customer information
//
// The application layer needs to look up customer names when creating orders
// but should NOT directly query the database or call external services itself.
//
// Currently customer data is in the same Postgres database, but this interface
// allows decoupling - customer data could move to a separate service later
// without changing any application or domain code.
type CustomerService interface {
	// GetCustomerName fetches the display name of a customer by ID
	// Returns empty string and error if customer is not found
	GetCustomerName(ctx context.Context, customerID string) (string, error)
}
