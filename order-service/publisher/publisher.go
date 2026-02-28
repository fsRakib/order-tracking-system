package publisher

import (
	"encoding/json"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

// OrderEvent is the message structure published to RabbitMQ
type OrderEvent struct {
	OrderID      string      `json:"order_id"`
	CustomerID   string      `json:"customer_id"`
	CustomerName string      `json:"customer_name"`
	Status       string      `json:"status"`
	TotalAmount  float64     `json:"total_amount"`
	Items        []OrderItem `json:"items"`
}

type OrderItem struct {
	SKU       string  `json:"sku"`
	Quantity  int32   `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

func Connect() {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatal("RABBITMQ_URL environment variable is not set")
	}

	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	log.Println("connected to RabbitMQ successfully")
}

func PublishOrderCreated(event OrderEvent) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("failed to open channel: %v", err)
		return
	}
	defer ch.Close()

	// Declare a fanout exchange so multiple consumers can receive the event
	err = ch.ExchangeDeclare(
		"order.events", // name
		"fanout",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Printf("failed to declare exchange: %v", err)
		return
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return
	}

	err = ch.Publish(
		"order.events", // exchange
		"",             // routing key (empty for fanout)
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("failed to publish event: %v", err)
		return
	}

	log.Printf("published OrderCreated event for order: %s", event.OrderID)
}
