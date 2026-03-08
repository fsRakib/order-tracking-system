package aggregate

import "order-service/domain/valueobject"

// OrderItem represents a single product line inside an Order
// It is NOT a standalone entity - it only exists as part of an Order aggregate
// Think of it like a row in an order table - meaningless without the order
type OrderItem struct {
	sku       valueobject.SKU
	quantity  int32
	unitPrice valueobject.Money
}

// newOrderItem is a package-private constructor
// Only the Order aggregate can create OrderItems (enforces encapsulation)
func newOrderItem(sku valueobject.SKU, quantity int32, unitPrice valueobject.Money) OrderItem {
	return OrderItem{
		sku:       sku,
		quantity:  quantity,
		unitPrice: unitPrice,
	}
}

// SKU returns the product SKU
func (i OrderItem) SKU() valueobject.SKU {
	return i.sku
}

// Quantity returns how many units were ordered
func (i OrderItem) Quantity() int32 {
	return i.quantity
}

// UnitPrice returns the price per unit
func (i OrderItem) UnitPrice() valueobject.Money {
	return i.unitPrice
}

// Subtotal calculates the total price for this line item
// Example: 3 units at $10.00 = $30.00
func (i OrderItem) Subtotal() valueobject.Money {
	return i.unitPrice.Multiply(i.quantity)
}
