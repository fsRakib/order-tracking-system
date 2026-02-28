package elastic

import (
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

var Client *elasticsearch.Client

func Connect() {
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		log.Fatal("ELASTICSEARCH_URL environment variable is not set")
	}

	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	var err error
	Client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create elasticsearch client: %v", err)
	}

	// Ping to verify connection
	res, err := Client.Info()
	if err != nil {
		log.Fatalf("failed to connect to elasticsearch: %v", err)
	}
	defer res.Body.Close()

	log.Println("connected to Elasticsearch successfully")

	// Create index with mapping if it doesn't exist
	createIndexIfNotExists()
}