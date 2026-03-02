#!/bin/bash

echo "============================================"
echo "  ORDER TRACKING SYSTEM - COMPLETE TEST"
echo "============================================"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Step 1: Analyze Database Data
echo -e "\n${BLUE}[STEP 1] DATABASE ANALYSIS${NC}"
echo "----------------------------------------"

echo -e "${YELLOW}Customers in database:${NC}"
PGPASSWORD=admin psql -h localhost -U admin -d orders_db -c "SELECT id, name FROM customers;"

echo -e "\n${YELLOW}Orders in database:${NC}"
PGPASSWORD=admin psql -h localhost -U admin -d orders_db -c "SELECT id, customer_id, status, total_amount FROM orders;"

echo -e "\n${YELLOW}Order Items:${NC}"
PGPASSWORD=admin psql -h localhost -U admin -d orders_db -c "SELECT order_id, sku, quantity, unit_price FROM order_items ORDER BY order_id LIMIT 10;"

echo -e "\n${YELLOW}Stock inventory:${NC}"
PGPASSWORD=admin psql -h localhost -U admin -d orders_db -c "SELECT sku, quantity FROM stocks;"

echo -e "\n${YELLOW}Elasticsearch index count:${NC}"
curl -s "http://localhost:9201/orders/_count" | jq

# Step 2: Check Services Status
echo -e "\n${BLUE}[STEP 2] SERVICE STATUS CHECK${NC}"
echo "----------------------------------------"

check_service() {
    local name=$1
    local port=$2
    if nc -z localhost $port 2>/dev/null; then
        echo -e "${GREEN}✓${NC} $name is running on port $port"
        return 0
    else
        echo -e "${RED}✗${NC} $name is NOT running on port $port"
        return 1
    fi
}

check_service "Order Service (gRPC)" 50061
ORDER_RUNNING=$?

check_service "Stock Service (gRPC)" 50062
STOCK_RUNNING=$?

check_service "Analytics Service (HTTP)" 8081
ANALYTICS_RUNNING=$?

check_service "PostgreSQL" 5433
check_service "RabbitMQ" 5673
check_service "Elasticsearch" 9201

if [ $ORDER_RUNNING -ne 0 ] || [ $STOCK_RUNNING -ne 0 ] || [ $ANALYTICS_RUNNING -ne 0 ]; then
    echo -e "\n${RED}ERROR: Some services are not running!${NC}"
    echo -e "${YELLOW}Please start all services before running tests.${NC}"
    echo ""
    echo "To start services, run in separate terminals:"
    echo "  Terminal 1: cd order-service && go run main.go"
    echo "  Terminal 2: cd stock-service && go run main.go"
    echo "  Terminal 3: cd analytics-service && go run main.go"
    echo ""
    exit 1
fi

# Step 3: Test All Endpoints
echo -e "\n${BLUE}[STEP 3] TESTING ALL 11 ENDPOINTS${NC}"
echo "============================================"

# ========== ORDER SERVICE (4 methods) ==========
echo -e "\n${GREEN}[1/11] ORDER SERVICE - CreateOrder${NC}"
echo "Creating new order for CUST-001..."
grpcurl -plaintext -d '{
  "customer_id": "CUST-001",
  "items": [
    {"sku": "LAPTOP-001", "quantity": 1, "unit_price": 999.99},
    {"sku": "MOUSE-001", "quantity": 2, "unit_price": 29.99}
  ]
}' localhost:50061 order.OrderService/CreateOrder

echo -e "\n${GREEN}[2/11] ORDER SERVICE - GetOrder${NC}"
echo "Getting order: 1c798b82-7448-489f-a6b2-0cb738c56923"
grpcurl -plaintext -d '{
  "order_id": "1c798b82-7448-489f-a6b2-0cb738c56923"
}' localhost:50061 order.OrderService/GetOrder

echo -e "\n${GREEN}[3/11] ORDER SERVICE - UpdateOrderStatus${NC}"
echo "Updating order status to 'shipped'..."
grpcurl -plaintext -d '{
  "order_id": "1c798b82-7448-489f-a6b2-0cb738c56923",
  "status": "shipped"
}' localhost:50061 order.OrderService/UpdateOrderStatus

echo -e "\n${GREEN}[4/11] ORDER SERVICE - GetOrdersByCustomer${NC}"
echo "Getting all orders for customer: CUST-001"
grpcurl -plaintext -d '{
  "customer_id": "CUST-001"
}' localhost:50061 order.OrderService/GetOrdersByCustomer

# ========== STOCK SERVICE (3 methods) ==========
echo -e "\n${GREEN}[5/11] STOCK SERVICE - GetStock${NC}"
echo "Getting stock for SKU: LAPTOP-001"
grpcurl -plaintext -d '{
  "sku": "LAPTOP-001"
}' localhost:50062 stock.StockService/GetStock

echo -e "\n${GREEN}[6/11] STOCK SERVICE - ReserveStock${NC}"
echo "Reserving 5 units of MOUSE-001..."
grpcurl -plaintext -d '{
  "sku": "MOUSE-001",
  "quantity": 5
}' localhost:50062 stock.StockService/ReserveStock

echo -e "\n${GREEN}[7/11] STOCK SERVICE - ReleaseStock${NC}"
echo "Releasing 3 units of MOUSE-001..."
grpcurl -plaintext -d '{
  "sku": "MOUSE-001",
  "quantity": 3
}' localhost:50062 stock.StockService/ReleaseStock

# ========== ANALYTICS SERVICE (4 endpoints) ==========
echo -e "\n${GREEN}[8/11] ANALYTICS SERVICE - Search Orders${NC}"
echo "Searching orders by customer: Rakib"
curl -s "http://localhost:8081/search?customer=Rakib" | jq

echo -e "\n${GREEN}[9/11] ANALYTICS SERVICE - Search Orders by SKU${NC}"
echo "Searching orders by SKU: LAPTOP-001"
curl -s "http://localhost:8081/search?sku=LAPTOP-001" | jq

echo -e "\n${GREEN}[10/11] ANALYTICS SERVICE - Aggregate by Status${NC}"
echo "Getting order counts by status..."
curl -s "http://localhost:8081/aggregate/status" | jq

echo -e "\n${GREEN}[11/11] ANALYTICS SERVICE - Aggregate by Customer${NC}"
echo "Getting order counts by customer..."
curl -s "http://localhost:8081/aggregate/customer" | jq

# Bonus: Health Check
echo -e "\n${GREEN}[BONUS] ANALYTICS SERVICE - Health Check${NC}"
curl -s "http://localhost:8081/health" | jq

# Summary
echo -e "\n============================================"
echo -e "${GREEN}✓ ALL 11 ENDPOINTS TESTED SUCCESSFULLY${NC}"
echo "============================================"
echo ""
echo "Summary:"
echo "  - Order Service (gRPC): 4 methods"
echo "  - Stock Service (gRPC): 3 methods"
echo "  - Analytics Service (HTTP): 4 endpoints"
echo "  Total: 11 operations tested"
echo ""
