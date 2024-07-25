package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/logger/config"
)

func main() {
	go config.ConsumeFromRabbitMQ()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3001"))
}
