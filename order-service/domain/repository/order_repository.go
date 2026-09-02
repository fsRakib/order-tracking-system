package repository

import (
	"context"

	"order-service/domain/aggregate"
	"order-service/domain/valueobject"
)

// OrderRepository defines how to persist and retrieve Order aggregates
//
// This is an INTERFACE - it lives in the domain layer but has NO implementation here
// The actual SQL/Postgres code is in infrastructure/persistence/
//
// Java equivalent:
//
//	public interface OrderRepository {
//	    void save(Order order);
//	    Optional<Order> findById(OrderID id);
//	}
//
// Why this matters: domain logic calls orderRepo.Save(order) without
// knowing whether the data goes to Postgres, MySQL, or even an in-memory map.
// Swapping databases = only change the infrastructure layer.
type OrderRepository interface {
	// Save persists a new Order to the database
	// Used when creating an order for the first time
	Save(ctx context.Context, order *aggregate.Order) error

	// Update persists changes to an existing Order
	// Used when order status changes or items are modified
	Update(ctx context.Context, order *aggregate.Order) error

	// FindByID retrieves a single Order by its unique ID
	// Returns nil and an error if the order does not exist
	FindByID(ctx context.Context, id valueobject.OrderID) (*aggregate.Order, error)

	// FindByCustomerID retrieves all Orders placed by a customer
	// Returns an empty slice (not error) if customer has no orders
	FindByCustomerID(ctx context.Context, customerID valueobject.CustomerID) ([]*aggregate.Order, error)
}
