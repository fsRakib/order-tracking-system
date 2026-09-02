package port

import (
	"context"

	"order-service/domain/event"
)

// EventPublisher defines how to publish domain events to a message broker
//
// This interface lives in the application layer because the application layer
// needs to publish events after saving an order, but should NOT know
// that the implementation uses RabbitMQ.
//
// The actual RabbitMQ implementation is in infrastructure/messaging/
//
// Java equivalent:
//
//	public interface EventPublisher {
//	    void publish(DomainEvent event);
//	}
type EventPublisher interface {
	// Publish sends a domain event to the message broker
	// The event type (OrderCreated, StatusUpdated, etc.) determines the routing
	Publish(ctx context.Context, domainEvent event.DomainEvent) error
}
