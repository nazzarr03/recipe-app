package config

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var (
	RabbitMQConn *amqp.Connection
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func ConnectRabbitMQ() {
	var err error
	RabbitMQConn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		panic(err)
	}
}

func PublishEmailToQueue(email, subject, body string) error {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(fmt.Sprintf("%s\n%s\n\n%s", email, subject, body)),
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		message,
	)
	if err != nil {
		fmt.Println("Failed to send email")
		return err
	}

	fmt.Println("Email sent to queue")
	return nil
}

func ConsumeEmailQueue() {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			parts := strings.SplitN(string(d.Body), "\n\n", 2)
			if len(parts) < 2 {
				log.Println("Invalid message format")
				d.Nack(false, false)
				continue
			}

			header := strings.SplitN(parts[0], "\n", 2)
			if len(header) < 2 {
				log.Println("Invalid header format")
				d.Nack(false, false)
				continue
			}

			email := header[0]
			subject := header[1]
			body := parts[1]

			from := os.Getenv("EMAIL_FROM")
			password := os.Getenv("EMAIL_PASSWORD")

			auth := smtp.PlainAuth("", from, password, os.Getenv("EMAIL_SMTP_SERVER"))
			to := []string{email}
			msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

			err := smtp.SendMail(fmt.Sprintf("%s:%s", os.Getenv("EMAIL_SMTP_SERVER"), os.Getenv("EMAIL_SMTP_PORT")), auth, from, to, msg)
			if err != nil {
				log.Printf("Failed to send email: %v", err)
				d.Nack(false, false)
			} else {
				log.Printf("Email sent to %s", email)
				d.Ack(false)
			}
		}
	}()
	<-forever
}
