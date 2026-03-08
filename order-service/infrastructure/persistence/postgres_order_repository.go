package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"order-service/domain/aggregate"
	"order-service/domain/repository"
	"order-service/domain/valueobject"
)

// PostgresOrderRepository is the concrete implementation of repository.OrderRepository
// It knows about SQL and Postgres - this is the ONLY place in the codebase that does
//
// Java equivalent: class PostgresOrderRepository implements OrderRepository { ... }
type PostgresOrderRepository struct {
	db *sql.DB
}

// NewPostgresOrderRepository creates a new repository instance
// Returns the interface type, not the concrete type - this enforces the contract
func NewPostgresOrderRepository(db *sql.DB) repository.OrderRepository {
	return &PostgresOrderRepository{db: db}
}

// Save persists a brand new Order aggregate to Postgres
// Converts Order aggregate fields back into database rows
func (r *PostgresOrderRepository) Save(ctx context.Context, order *aggregate.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Insert the order row
	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders (id, customer_id, total_amount, status) VALUES ($1, $2, $3, $4)`,
		order.ID().Value(),
		order.CustomerID().Value(),
		order.TotalAmount().Amount(),
		order.Status().String(),
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Insert each order item row
	for _, item := range order.Items() {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO order_items (order_id, sku, quantity, unit_price) VALUES ($1, $2, $3, $4)`,
			order.ID().Value(),
			item.SKU().Value(),
			item.Quantity(),
			item.UnitPrice().Amount(),
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	return tx.Commit()
}

// Update persists changes to an existing Order (e.g., status change)
func (r *PostgresOrderRepository) Update(ctx context.Context, order *aggregate.Order) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status = $1, total_amount = $2 WHERE id = $3`,
		order.Status().String(),
		order.TotalAmount().Amount(),
		order.ID().Value(),
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// FindByID retrieves an Order aggregate from Postgres by its ID
// Reconstructs the full aggregate including all items
func (r *PostgresOrderRepository) FindByID(ctx context.Context, id valueobject.OrderID) (*aggregate.Order, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT o.id, o.customer_id, c.name, o.total_amount, o.status
		 FROM orders o
		 LEFT JOIN customers c ON c.id = o.customer_id
		 WHERE o.id = $1`,
		id.Value(),
	)

	return r.scanOrder(ctx, row)
}

// FindByCustomerID retrieves all Orders for a given customer
func (r *PostgresOrderRepository) FindByCustomerID(ctx context.Context, customerID valueobject.CustomerID) ([]*aggregate.Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT o.id, o.customer_id, c.name, o.total_amount, o.status
		 FROM orders o
		 LEFT JOIN customers c ON c.id = o.customer_id
		 WHERE o.customer_id = $1
		 ORDER BY o.created_at DESC`,
		customerID.Value(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []*aggregate.Order
	for rows.Next() {
		order, err := r.scanOrder(ctx, rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// ---- Private Helpers ----

// rowScanner is an interface shared by *sql.Row and *sql.Rows
// This allows scanOrder to work with both QueryRow and Query results
type rowScanner interface {
	Scan(dest ...any) error
}

// scanOrder reads a database row and reconstructs an Order aggregate
// This is the reverse of Save - DB rows → domain aggregate
func (r *PostgresOrderRepository) scanOrder(ctx context.Context, row rowScanner) (*aggregate.Order, error) {
	var (
		rawID           string
		rawCustomerID   string
		rawCustomerName sql.NullString
		rawTotal        float64
		rawStatus       string
	)

	err := row.Scan(&rawID, &rawCustomerID, &rawCustomerName, &rawTotal, &rawStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to scan order row: %w", err)
	}

	// Reconstruct value objects from raw DB strings
	orderID, err := valueobject.NewOrderIDFromString(rawID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID in database: %w", err)
	}

	customerID, err := valueobject.NewCustomerID(rawCustomerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID in database: %w", err)
	}

	status, err := valueobject.NewOrderStatus(rawStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid order status in database: %w", err)
	}

	totalAmount, err := valueobject.NewMoney(rawTotal, "USD")
	if err != nil {
		return nil, fmt.Errorf("invalid total amount in database: %w", err)
	}

	customerName := ""
	if rawCustomerName.Valid {
		customerName = rawCustomerName.String
	}

	// Fetch the order items from the order_items table
	items, err := r.fetchOrderItems(ctx, orderID.Value())
	if err != nil {
		return nil, err
	}

	// Use ReconstructOrder (not NewOrder) because this is NOT a new order
	return aggregate.ReconstructOrder(orderID, customerID, customerName, items, status, totalAmount), nil
}

// fetchOrderItems queries the order_items table for a given order ID
func (r *PostgresOrderRepository) fetchOrderItems(ctx context.Context, orderID string) ([]aggregate.OrderItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT sku, quantity, unit_price FROM order_items WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}
	defer rows.Close()

	var items []aggregate.OrderItem
	for rows.Next() {
		var rawSKU string
		var quantity int32
		var rawUnitPrice float64

		if err := rows.Scan(&rawSKU, &quantity, &rawUnitPrice); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		sku, err := valueobject.NewSKU(rawSKU)
		if err != nil {
			return nil, fmt.Errorf("invalid SKU in database: %w", err)
		}

		unitPrice, err := valueobject.NewMoney(rawUnitPrice, "USD")
		if err != nil {
			return nil, fmt.Errorf("invalid unit price in database: %w", err)
		}

		items = append(items, aggregate.NewOrderItem(sku, quantity, unitPrice))
	}

	return items, nil
}
