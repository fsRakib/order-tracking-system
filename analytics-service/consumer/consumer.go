package consumer

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"analytics-service/elastic"

	amqp "github.com/rabbitmq/amqp091-go"
)

// IncomingOrderEvent matches the structure published by the Order Service
type IncomingOrderEvent struct {
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

	// Declare the same fanout exchange the Order Service publishes to
	err = ch.ExchangeDeclare(
		"order.events",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}

	// Create a queue and bind it to the exchange
	q, err := ch.QueueDeclare(
		"analytics.order.events", // unique queue name for this consumer
		true,                     // durable
		false,                    // auto-delete
		false,                    // exclusive
		false,                    // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	// Bind the queue to the fanout exchange
	err = ch.QueueBind(q.Name, "", "order.events", false, nil)
	if err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to register consumer: %v", err)
	}

	log.Println("analytics consumer listening on order.events exchange...")

	for msg := range msgs {
		var event IncomingOrderEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("failed to parse message: %v", err)
			continue
		}

		log.Printf("received order event: %s", event.OrderID)

		// Convert to Elasticsearch document
		var esItems []elastic.OrderItem
		for _, item := range event.Items {
			esItems = append(esItems, elastic.OrderItem{
				SKU:       item.SKU,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			})
		}

		doc := elastic.OrderDocument{
			OrderID:      event.OrderID,
			CustomerID:   event.CustomerID,
			CustomerName: event.CustomerName,
			Status:       event.Status,
			TotalAmount:  event.TotalAmount,
			Items:        esItems,
			CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		}

		if err := elastic.IndexOrder(doc); err != nil {
			log.Printf("failed to index order %s: %v", event.OrderID, err)
		}
	}
}
