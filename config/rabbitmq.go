package config

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

var (
	RabbitMQConn *amqp.Connection
)

func ConnectRabbitMQ() *amqp.Connection {
	var err error
	RabbitMQConn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		panic("failed to connect to RabbitMQ")
	}

	fmt.Println("RabbitMQ connected successfully!")

	return RabbitMQConn
}
