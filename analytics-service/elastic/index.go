package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const indexName = "orders"

// OrderDocument is the structure stored in Elasticsearch
type OrderDocument struct {
	OrderID      string      `json:"order_id"`
	CustomerID   string      `json:"customer_id"`
	CustomerName string      `json:"customer_name"`
	Status       string      `json:"status"`
	TotalAmount  float64     `json:"total_amount"`
	Items        []OrderItem `json:"items"`
	CreatedAt    string      `json:"created_at"`
}

type OrderItem struct {
	SKU       string  `json:"sku"`
	Quantity  int32   `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

// createIndexIfNotExists creates the orders index with proper field mappings.
// Returns true if the index was freshly created (caller should trigger a sync).
func createIndexIfNotExists() bool {
	mapping := `{
		"mappings": {
			"properties": {
				"order_id":      { "type": "keyword" },
				"customer_id":   { "type": "keyword" },
				"customer_name": { 
					"type": "text",
					"fields": {
						"keyword": { "type": "keyword" }
					}
				},
				"status":        { "type": "keyword" },
				"total_amount":  { "type": "float" },
				"created_at":    { "type": "date" },
				"items": {
					"type": "nested",
					"properties": {
						"sku":        { "type": "keyword" },
						"quantity":   { "type": "integer" },
						"unit_price": { "type": "float" }
					}
				}
			}
		}
	}`

	res, err := Client.Indices.Exists([]string{indexName})
	if err != nil {
		log.Fatalf("failed to check index existence: %v", err)
	}
	defer res.Body.Close()

	// 404 means index does not exist — create it
	if res.StatusCode == 404 {
		createRes, err := Client.Indices.Create(
			indexName,
			Client.Indices.Create.WithBody(strings.NewReader(mapping)),
		)
		if err != nil {
			log.Fatalf("failed to create index: %v", err)
		}
		defer createRes.Body.Close()
		log.Printf("created Elasticsearch index: %s", indexName)
		return true // freshly created — caller should sync historical data
	}

	log.Printf("Elasticsearch index already exists: %s", indexName)
	return false
}

// IndexOrder stores an order document in Elasticsearch
func IndexOrder(doc OrderDocument) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %v", err)
	}

	res, err := Client.Index(
		indexName,
		bytes.NewReader(body),
		Client.Index.WithDocumentID(doc.OrderID),
		Client.Index.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("failed to index document: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch index error: %s", res.String())
	}

	log.Printf("indexed order: %s", doc.OrderID)
	return nil
}

// SearchOrders searches by customer name, product SKU, or status
func SearchOrders(customerName, sku, status string) ([]OrderDocument, error) {
	var mustClauses []map[string]interface{}

	if customerName != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"wildcard": map[string]interface{}{
				"customer_name.keyword": map[string]interface{}{
					"value":            "*" + customerName + "*",
					"case_insensitive": true,
				},
			},
		})
	}

	if status != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"status": status,
			},
		})
	}

	if sku != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "items",
				"query": map[string]interface{}{
					"term": map[string]interface{}{
						"items.sku": sku,
					},
				},
			},
		})
	}

	// If no filters provided, match all
	var query map[string]interface{}
	if len(mustClauses) == 0 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": mustClauses,
				},
			},
		}
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %v", err)
	}

	res, err := Client.Search(
		Client.Search.WithIndex(indexName),
		Client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("search error: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch search error: %s", res.String())
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return parseHits(result)
}

// AggregateByStatus returns order counts grouped by status
func AggregateByStatus() (map[string]int, error) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"orders_by_status": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "status",
				},
			},
		},
	}

	return runAggregation(query, "orders_by_status")
}

// AggregateByCustomer returns order counts grouped by customer
func AggregateByCustomer() (map[string]int, error) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"orders_by_customer": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "customer_id",
				},
			},
		},
	}

	return runAggregation(query, "orders_by_customer")
}

// runAggregation is a helper that executes an aggregation query
func runAggregation(query map[string]interface{}, aggName string) (map[string]int, error) {
	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal aggregation query: %v", err)
	}

	res, err := Client.Search(
		Client.Search.WithIndex(indexName),
		Client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("aggregation error: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch aggregation error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	counts := make(map[string]int)
	aggs := result["aggregations"].(map[string]interface{})
	buckets := aggs[aggName].(map[string]interface{})["buckets"].([]interface{})
	for _, b := range buckets {
		bucket := b.(map[string]interface{})
		key := bucket["key"].(string)
		count := int(bucket["doc_count"].(float64))
		counts[key] = count
	}

	return counts, nil
}

// parseHits extracts OrderDocuments from an ES search response
func parseHits(result map[string]interface{}) ([]OrderDocument, error) {
	var orders []OrderDocument

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		data, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}
		var doc OrderDocument
		if err := json.Unmarshal(data, &doc); err != nil {
			return nil, err
		}
		orders = append(orders, doc)
	}

	return orders, nil
}
