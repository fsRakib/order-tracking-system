package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"analytics-service/elastic"

	_ "github.com/lib/pq"
)

// SyncOrdersToElasticsearch queries all orders from PostgreSQL and indexes
// them into Elasticsearch. Called on startup when a fresh index is created.
func SyncOrdersToElasticsearch() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Println("DB_URL not set, skipping historical sync")
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("sync: failed to open DB connection: %v", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("sync: failed to ping DB: %v", err)
		return
	}

	log.Println("sync: starting historical order sync from PostgreSQL...")

	rows, err := db.Query(`
		SELECT o.id, o.customer_id, c.name, o.status, o.total_amount, o.created_at
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		ORDER BY o.created_at ASC
	`)
	if err != nil {
		log.Printf("sync: failed to query orders: %v", err)
		return
	}
	defer rows.Close()

	var synced int
	for rows.Next() {
		var doc elastic.OrderDocument
		var createdAt time.Time

		if err := rows.Scan(
			&doc.OrderID, &doc.CustomerID, &doc.CustomerName,
			&doc.Status, &doc.TotalAmount, &createdAt,
		); err != nil {
			log.Printf("sync: failed to scan order row: %v", err)
			continue
		}
		doc.CreatedAt = createdAt.UTC().Format(time.RFC3339)

		// Fetch items for this order
		items, err := fetchItems(db, doc.OrderID)
		if err != nil {
			log.Printf("sync: failed to fetch items for order %s: %v", doc.OrderID, err)
			continue
		}
		doc.Items = items

		if err := elastic.IndexOrder(doc); err != nil {
			log.Printf("sync: failed to index order %s: %v", doc.OrderID, err)
			continue
		}
		synced++
	}

	log.Printf("sync: completed — indexed %d orders from PostgreSQL", synced)
}

func fetchItems(db *sql.DB, orderID string) ([]elastic.OrderItem, error) {
	rows, err := db.Query(
		`SELECT sku, quantity, unit_price FROM order_items WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}
	defer rows.Close()

	var items []elastic.OrderItem
	for rows.Next() {
		var item elastic.OrderItem
		if err := rows.Scan(&item.SKU, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}
