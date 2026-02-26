package consumer

import (
	"encoding/json"
	"log"
	"os"

	"stock-service/db"

	amqp "github.com/rabbitmq/amqp091-go"
)

// CancelledOrderEvent is the message structure we expect from order.cancelled queue
type CancelledOrderEvent struct {
	SKU      string `json:"sku"`
	Quantity int32  `json:"quantity"`
}

func StartConsumer() {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatal("RABBITMQ_URL environment variable is not set")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue (safe to call even if it already exists)
	_, err = ch.QueueDeclare(
		"order.cancelled", // queue name
		true,              // durable
		false,             // auto-delete
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		"order.cancelled", // queue
		"",                // consumer tag
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		log.Fatalf("failed to register consumer: %v", err)
	}

	log.Println("stock consumer listening on order.cancelled queue...")

	for msg := range msgs {
		var event CancelledOrderEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("failed to parse message: %v", err)
			continue
		}

		log.Printf("received cancellation: SKU=%s, Quantity=%d", event.SKU, event.Quantity)

		if err := db.ReleaseStock(event.SKU, event.Quantity); err != nil {
			log.Printf("failed to release stock for SKU %s: %v", event.SKU, err)
		} else {
			log.Printf("stock restored for SKU %s, quantity %d", event.SKU, event.Quantity)
		}
	}
}
