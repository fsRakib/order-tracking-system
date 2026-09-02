# Domain Layer

This layer contains the core business logic and rules. It has no dependencies on external frameworks or infrastructure.

## Structure

- **aggregate/**: Aggregate roots (Order) - entities with identity and lifecycle
- **valueobject/**: Value objects (Money, OrderStatus) - immutable objects defined by their attributes
- **repository/**: Repository interfaces (no implementation here)
- **event/**: Domain events that represent business occurrences
