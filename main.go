package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nazzarr03/recipe-app/controllers"
	"github.com/nazzarr03/recipe-app/middleware"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	app := fiber.New()

	secret := os.Getenv("JWT_SECRET")
	app.Use(middleware.Authentication(secret))

	app.Post("/signup", controllers.SignUp)
	app.Post("/login", controllers.Login)

	app.Listen(":3000")
}
