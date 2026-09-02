package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"order-service/application/port"
	"order-service/domain/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

const exchangeName = "order.events"

// RabbitMQPublisher implements port.EventPublisher using RabbitMQ
// This is the ONLY place in the codebase that knows about RabbitMQ
//
// Java equivalent: class RabbitMQPublisher implements EventPublisher { ... }
type RabbitMQPublisher struct {
	conn *amqp.Connection
}

// NewRabbitMQPublisher creates a new publisher connected to RabbitMQ
// Returns the interface type to enforce the contract
func NewRabbitMQPublisher(conn *amqp.Connection) port.EventPublisher {
	return &RabbitMQPublisher{conn: conn}
}

// Publish sends a domain event to the RabbitMQ fanout exchange
// Automatically serializes different event types to JSON
func (p *RabbitMQPublisher) Publish(ctx context.Context, domainEvent event.DomainEvent) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	// Declare the fanout exchange (safe to declare multiple times)
	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,  // durable - survives broker restart
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Serialize the event to JSON bytes
	body, err := json.Marshal(domainEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish to the fanout exchange with the event type as routing key
	err = ch.Publish(
		exchangeName,
		domainEvent.EventType(), // routing key (fanout ignores this but useful for tracing)
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("published event [%s] for order", domainEvent.EventType())
	return nil
}

// Connect establishes a connection to RabbitMQ and returns it
// Called once during application startup in main.go
func Connect(url string) (*amqp.Connection, error) {
	if url == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is not set")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	log.Println("connected to RabbitMQ successfully")
	return conn, nil
}
