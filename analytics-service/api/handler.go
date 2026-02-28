package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"analytics-service/elastic"
)

func StartHTTPServer() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Search orders by customer name, sku, or status
	// GET /search?customer=John&sku=SKU-001&status=confirmed
	mux.HandleFunc("/search", searchHandler)

	// Aggregation by status
	// GET /aggregate/status
	mux.HandleFunc("/aggregate/status", aggregateByStatusHandler)

	// Aggregation by customer
	// GET /aggregate/customer
	mux.HandleFunc("/aggregate/customer", aggregateByCustomerHandler)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("analytics HTTP server running on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customerName := r.URL.Query().Get("customer")
	sku := r.URL.Query().Get("sku")
	status := r.URL.Query().Get("status")

	orders, err := elastic.SearchOrders(customerName, sku, status)
	if err != nil {
		log.Printf("search error: %v", err)
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"total":  len(orders),
		"orders": orders,
	})
}

func aggregateByStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	counts, err := elastic.AggregateByStatus()
	if err != nil {
		log.Printf("aggregation error: %v", err)
		http.Error(w, "aggregation failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"aggregation": "orders_by_status",
		"data":        counts,
	})
}

func aggregateByCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	counts, err := elastic.AggregateByCustomer()
	if err != nil {
		log.Printf("aggregation error: %v", err)
		http.Error(w, "aggregation failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"aggregation": "orders_by_customer",
		"data":        counts,
	})
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}