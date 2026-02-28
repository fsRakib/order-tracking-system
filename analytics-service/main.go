package main

import (
	"log"

	"analytics-service/api"
	"analytics-service/consumer"
	"analytics-service/elastic"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	// Connect to Elasticsearch
	elastic.Connect()

	// Start RabbitMQ consumer in background
	go consumer.StartConsumer()

	// Start HTTP API server (blocking)
	api.StartHTTPServer()
}