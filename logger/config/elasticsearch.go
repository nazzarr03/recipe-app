package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/nazzarr03/logger/models"
)

func ConnectElasticsearch() {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	var err error
	EsClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esURL},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	log.Println("Connected to Elasticsearch")
}

func SendLogToElasticsearch(logMessage models.LogMessage) error {
	log.Printf("Received log elastic: %+v\n", logMessage)
	jsonLog, err := json.Marshal(logMessage)
	if err != nil {
		return fmt.Errorf("error marshalling the log: %w", err)
	}

	log.Printf("Sending log message to Elasticsearch: %s\n", string(jsonLog))

	req := esapi.IndexRequest{
		Index:   "logs",
		Body:    strings.NewReader(string(jsonLog)),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), EsClient)
	if err != nil {
		return fmt.Errorf("error indexing the log: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	log.Println("Log message sent to Elasticsearch")
	return nil
}
