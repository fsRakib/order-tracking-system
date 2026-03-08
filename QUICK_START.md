# Quick Start Guide

> This guide uses Docker Compose V2 (`docker compose`). Install Docker Desktop or Docker Engine with Compose V2 plugin if not already available.

## One-Command Setup

```bash
# Navigate to project directory
cd order-tracking-system

# Start everything with Docker
docker compose up --build -d
```

All 6 services start automatically.

---

## Verify Everything is Running

```bash
docker compose ps
```

Expected containers running:

- order-tracking-postgres
- order-tracking-rabbitmq
- order-tracking-elasticsearch
- order-tracking-order-service
- order-tracking-stock-service
- order-tracking-analytics-service

---

## Insert Sample Data

```bash
docker exec order-tracking-postgres psql -U admin -d orders_db -c "
INSERT INTO customers (id, name) VALUES
('CUST-001', 'Md. Rakibul Kabir'),
('CUST-002', 'Arka Das'),
('CUST-003', 'Morshed Alam'),
('CUST-004', 'Emran Ahmed Emon')
ON CONFLICT (id) DO NOTHING;

INSERT INTO stocks (sku, quantity) VALUES
('LAPTOP-001', 50),
('MOUSE-001', 100),
('KEYBOARD-001', 75),
('MONITOR-001', 30),
('HEADSET-001', 60)
ON CONFLICT (sku) DO NOTHING;"
```

---

## Test the System

### 🔹 ORDER SERVICE (gRPC) - 4 Methods

**1. CreateOrder**

```bash
grpcurl -plaintext -d '{
  "customer_id": "CUST-001",
  "items": [
    {"sku": "LAPTOP-001", "quantity": 1, "unit_price": 999.99}
  ]
}' localhost:50061 order.OrderService/CreateOrder
```

**2. GetOrder** (replace ORDER_ID with the one from CreateOrder response)

```bash
grpcurl -plaintext -d '{
  "order_id": "YOUR_ORDER_ID_HERE"
}' localhost:50061 order.OrderService/GetOrder
```

**3. UpdateOrderStatus**

```bash
grpcurl -plaintext -d '{
  "order_id": "YOUR_ORDER_ID_HERE",
  "status": "shipped"
}' localhost:50061 order.OrderService/UpdateOrderStatus
```

**4. GetOrdersByCustomer**

```bash
grpcurl -plaintext -d '{
  "customer_id": "CUST-001"
}' localhost:50061 order.OrderService/GetOrdersByCustomer
```

---

### 🔹 STOCK SERVICE (gRPC) - 3 Methods

**5. GetStock**

```bash
grpcurl -plaintext -d '{
  "sku": "LAPTOP-001"
}' localhost:50062 stock.StockService/GetStock
```

**6. ReserveStock**

```bash
grpcurl -plaintext -d '{
  "sku": "MOUSE-001",
  "quantity": 5
}' localhost:50062 stock.StockService/ReserveStock
```

**7. ReleaseStock**

```bash
grpcurl -plaintext -d '{
  "sku": "MOUSE-001",
  "quantity": 3
}' localhost:50062 stock.StockService/ReleaseStock
```

---

### 🔹 ANALYTICS SERVICE (HTTP REST) - 4 Endpoints

**8. Search by Customer Name**

```bash
curl -s "http://localhost:8081/search?customer=Rakib" | jq
```

**9. Search by SKU**

```bash
curl -s "http://localhost:8081/search?sku=LAPTOP-001" | jq
```

**10. Aggregate by Status**

```bash
curl -s "http://localhost:8081/aggregate/status" | jq
```

**11. Aggregate by Customer**

```bash
curl -s "http://localhost:8081/aggregate/customer" | jq
```

**Bonus: Health Check**

```bash
curl -s "http://localhost:8081/health" | jq
```

---

### 📊 Complete Test Flow

Run all tests in sequence:

```bash
# 1. Create order and capture ID
ORDER_ID=$(grpcurl -plaintext -d '{"customer_id":"CUST-001","items":[{"sku":"LAPTOP-001","quantity":1,"unit_price":999.99}]}' localhost:50061 order.OrderService/CreateOrder | grep -o '"orderId":"[^"]*"' | cut -d'"' -f4)

echo "Created Order: $ORDER_ID"

# 2. Get order details
grpcurl -plaintext -d "{\"order_id\":\"$ORDER_ID\"}" localhost:50061 order.OrderService/GetOrder

# 3. Check stock
grpcurl -plaintext -d '{"sku":"LAPTOP-001"}' localhost:50062 stock.StockService/GetStock

# 4. Update status
grpcurl -plaintext -d "{\"order_id\":\"$ORDER_ID\",\"status\":\"shipped\"}" localhost:50061 order.OrderService/UpdateOrderStatus

# 5. Wait for Elasticsearch indexing
sleep 3

# 6. Search in analytics
curl -s "http://localhost:8081/search?customer=Rakib" | jq

# 7. View aggregations
curl -s "http://localhost:8081/aggregate/status" | jq
```

---

## Access Services

- **Order Service:** localhost:50061 (gRPC)
- **Stock Service:** localhost:50062 (gRPC)
- **Analytics API:** http://localhost:8081
- **RabbitMQ UI:** http://localhost:15673 (admin/admin)
- **Elasticsearch:** http://localhost:9201

---

## Stop Everything

```bash
docker compose down
```

---

**For full documentation, architecture details, and DDD design explanation, see README.md.**
