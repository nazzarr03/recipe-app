package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"
	"github.com/nazzarr03/logger/models"
	"github.com/streadway/amqp"
)

var (
	EsClient     *elasticsearch.Client
	RabbitMQConn *amqp.Connection
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	ConnectRabbitMQ()
	ConnectElasticsearch()
}

func ConnectRabbitMQ() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	var err error
	RabbitMQConn, err = amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	log.Println("Connected to RabbitMQ")
}

func ConsumeFromRabbitMQ() {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Failed to get RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"logger",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare RabbitMQ queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	go func() {
		for d := range msgs {
			logMessage := models.LogMessage{}
			messageBody := string(d.Body)
			fmt.Println("Received message body: ", messageBody)

			if err := json.Unmarshal([]byte(messageBody), &logMessage); err != nil {
				log.Printf("Failed to unmarshal log message: %v", err)
				continue
			}

			log.Printf("Processing logMessage: %v", logMessage)
			if err := SendLogToElasticsearch(logMessage); err != nil {
				log.Printf("Failed to send log message to Elasticsearch: %v", err)
			}
		}
	}()

	for {
		time.Sleep(time.Second)
	}
}
