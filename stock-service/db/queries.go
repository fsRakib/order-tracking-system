package db

import (
	"database/sql"
	"fmt"
)

type Stock struct {
	ID       int
	SKU      string
	Quantity int
}

// GetStockBySKU fetches a single stock row by SKU
func GetStockBySKU(sku string) (*Stock, error) {
	row := DB.QueryRow(`SELECT id, sku, quantity FROM stocks WHERE sku = $1`, sku)

	var s Stock
	if err := row.Scan(&s.ID, &s.SKU, &s.Quantity); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("SKU not found: %s", sku)
		}
		return nil, fmt.Errorf("query error: %v", err)
	}
	return &s, nil
}

// ReserveStock decreases stock quantity inside a transaction using SELECT FOR UPDATE
func ReserveStock(sku string, quantity int32) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	var available int32
	err = tx.QueryRow(
		`SELECT quantity FROM stocks WHERE sku = $1 FOR UPDATE`, sku,
	).Scan(&available)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("SKU not found: %s", sku)
	}

	if available < quantity {
		tx.Rollback()
		return fmt.Errorf("insufficient stock for SKU %s: available %d, requested %d", sku, available, quantity)
	}

	_, err = tx.Exec(
		`UPDATE stocks SET quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP WHERE sku = $2`,
		quantity, sku,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to reserve stock: %v", err)
	}

	return tx.Commit()
}

// ReleaseStock increases stock quantity inside a transaction
func ReleaseStock(sku string, quantity int32) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(
		`UPDATE stocks SET quantity = quantity + $1, updated_at = CURRENT_TIMESTAMP WHERE sku = $2`,
		quantity, sku,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to release stock: %v", err)
	}

	return tx.Commit()
}