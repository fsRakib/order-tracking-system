package db

import (
	"database/sql"
	"fmt"
)

type Order struct {
	ID          string
	CustomerID  string
	TotalAmount float64
	Status      string
	CreatedAt   string
	Items       []OrderItem
}

type OrderItem struct {
	SKU       string
	Quantity  int32
	UnitPrice float64
}

// InsertOrder saves a new order and its items inside a single transaction
func InsertOrder(order Order) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(
		`INSERT INTO orders (id, customer_id, total_amount, status) VALUES ($1, $2, $3, $4)`,
		order.ID, order.CustomerID, order.TotalAmount, order.Status,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert order: %v", err)
	}

	for _, item := range order.Items {
		_, err = tx.Exec(
			`INSERT INTO order_items (order_id, sku, quantity, unit_price) VALUES ($1, $2, $3, $4)`,
			order.ID, item.SKU, item.Quantity, item.UnitPrice,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert order item: %v", err)
		}
	}

	return tx.Commit()
}

// GetOrderByID fetches a single order with its items
func GetOrderByID(orderID string) (*Order, error) {
	row := DB.QueryRow(
		`SELECT id, customer_id, total_amount, status, created_at FROM orders WHERE id = $1`,
		orderID,
	)

	var o Order
	if err := row.Scan(&o.ID, &o.CustomerID, &o.TotalAmount, &o.Status, &o.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found: %s", orderID)
		}
		return nil, fmt.Errorf("query error: %v", err)
	}

	items, err := getOrderItems(orderID)
	if err != nil {
		return nil, err
	}
	o.Items = items

	return &o, nil
}

// GetOrdersByCustomerID fetches all orders for a customer
func GetOrdersByCustomerID(customerID string) ([]Order, error) {
	rows, err := DB.Query(
		`SELECT id, customer_id, total_amount, status, created_at FROM orders WHERE customer_id = $1`,
		customerID,
	)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.TotalAmount, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}

		items, err := getOrderItems(o.ID)
		if err != nil {
			return nil, err
		}
		o.Items = items
		orders = append(orders, o)
	}

	return orders, nil
}

// UpdateOrderStatus changes the status of an order
func UpdateOrderStatus(orderID, status string) (*Order, error) {
	_, err := DB.Exec(
		`UPDATE orders SET status = $1 WHERE id = $2`,
		status, orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %v", err)
	}

	return GetOrderByID(orderID)
}

// getOrderItems is a helper to fetch items for a given order
func getOrderItems(orderID string) ([]OrderItem, error) {
	rows, err := DB.Query(
		`SELECT sku, quantity, unit_price FROM order_items WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %v", err)
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.SKU, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}