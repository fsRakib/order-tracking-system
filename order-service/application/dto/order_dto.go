package dto

// OrderItemDTO carries item data across layer boundaries
// DTOs are plain data structs - no business logic, no methods
// Think of them as the "shape" of data that crosses into/out of the application layer
type OrderItemDTO struct {
	SKU       string  `json:"sku"`
	Quantity  int32   `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
}

// OrderDTO carries order data from application layer to the gRPC handler
// The gRPC handler converts this into a protobuf response
type OrderDTO struct {
	OrderID      string         `json:"order_id"`
	CustomerID   string         `json:"customer_id"`
	CustomerName string         `json:"customer_name"`
	Status       string         `json:"status"`
	TotalAmount  float64        `json:"total_amount"`
	Items        []OrderItemDTO `json:"items"`
}
