package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/recipeApp/config"
	"github.com/streadway/amqp"
)

type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Latency   string `json:"latency"`
	Method    string `json:"method"`
	Path      string `json:"path"`
}

func SendLogToLogger(logMessage LogMessage) error {
	ch, err := config.RabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
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
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	messageBody, err := json.Marshal(logMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal log message: %w", err)
	}

	log.Println(string(messageBody))

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        messageBody,
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		message,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Log message sent to logger: %v", logMessage)

	return nil
}

func LogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		logMessage := LogMessage{
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    c.Response().StatusCode(),
			Latency:   time.Since(start).String(),
			Method:    c.Method(),
			Path:      c.Path(),
		}

		go func() {
			if err := SendLogToLogger(logMessage); err != nil {
				fmt.Printf("Error sending log to logger: %v\n", err)
			}
		}()

		return err
	}
}
