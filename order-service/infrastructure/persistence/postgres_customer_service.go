package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"order-service/application/port"
)

// PostgresCustomerService implements port.CustomerService
// Fetches customer data from the same Postgres database
type PostgresCustomerService struct {
	db *sql.DB
}

// NewPostgresCustomerService creates a new CustomerService backed by Postgres
func NewPostgresCustomerService(db *sql.DB) port.CustomerService {
	return &PostgresCustomerService{db: db}
}

// GetCustomerName fetches the display name of a customer by their ID
func (s *PostgresCustomerService) GetCustomerName(ctx context.Context, customerID string) (string, error) {
	var name string

	err := s.db.QueryRowContext(ctx,
		`SELECT name FROM customers WHERE id = $1`,
		customerID,
	).Scan(&name)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("customer not found: %s", customerID)
		}
		return "", fmt.Errorf("failed to fetch customer name: %w", err)
	}

	return name, nil
}
