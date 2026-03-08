package main

import (
	"log"

	"analytics-service/api"
	"analytics-service/consumer"
	"analytics-service/db"
	"analytics-service/elastic"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	// Connect to Elasticsearch; freshIndex = true means the index was just created
	freshIndex := elastic.Connect()

	// If the index is brand new, backfill all historical orders from PostgreSQL
	if freshIndex {
		log.Println("fresh index detected — syncing historical orders from PostgreSQL")
		db.SyncOrdersToElasticsearch()
	}

	// Start RabbitMQ consumer in background
	go consumer.StartConsumer()

	// Start HTTP API server (blocking)
	api.StartHTTPServer()
}