package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nazzarr03/mailService/config"
)

func main() {
	config.ConnectRabbitMQ()
	defer config.RabbitMQConn.Close()

	go config.ConsumeEmailQueue()

	fmt.Println("Mail service started")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("Shutting down mail service...")
}
