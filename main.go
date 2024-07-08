package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/recipe-app/controllers"
	"github.com/nazzarr03/recipe-app/middleware"
	"github.com/nazzarr03/recipe-app/routes"
)

func main() {
	app := fiber.New()

	app.Post("/signup", controllers.SignUp)
	app.Post("/login", controllers.Login)

	app.Use(middleware.Authentication())

	routes.RecipeRoutes(app)

	app.Listen(":3000")
}
