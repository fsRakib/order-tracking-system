#!/bin/bash
echo "=== Testing Order Service ==="
echo "1. CreateOrder"
grpcurl -plaintext -d '{"customer_id":"CUST-999","items":[{"sku":"LAPTOP-001","quantity":1,"unit_price":999.99}]}' localhost:50051 order.OrderService/CreateOrder

echo -e "\n2. GetOrder"
grpcurl -plaintext -d '{"order_id":"1c798b82-7448-489f-a6b2-0cb738c56923"}' localhost:50051 order.OrderService/GetOrder

echo -e "\n3. GetOrdersByCustomer"
grpcurl -plaintext -d '{"customer_id":"CUST-001"}' localhost:50051 order.OrderService/GetOrdersByCustomer

echo -e "\n=== Testing Stock Service ==="
echo "4. GetStock"
grpcurl -plaintext -d '{"sku":"LAPTOP-001"}' localhost:50052 stock.StockService/GetStock

echo -e "\n5. ReserveStock"
grpcurl -plaintext -d '{"sku":"MOUSE-001","quantity":5}' localhost:50052 stock.StockService/ReserveStock

echo -e "\n6. ReleaseStock"
grpcurl -plaintext -d '{"sku":"MOUSE-001","quantity":3}' localhost:50052 stock.StockService/ReleaseStock

echo -e "\n=== Testing Analytics Service ==="
echo "7. Search Orders"
curl -s "http://localhost:8080/search?customer=Rakib" | jq

echo -e "\n8. Aggregate by Status"
curl -s "http://localhost:8080/aggregate/status" | jq

echo -e "\n9. Aggregate by Customer"
curl -s "http://localhost:8080/aggregate/customer" | jq

echo -e "\n10. Health Check"
curl -s "http://localhost:8080/health" | jq

echo -e "\n=== All Tests Complete ==="
