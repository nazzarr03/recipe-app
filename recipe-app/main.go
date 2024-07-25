package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/recipeApp/controllers"
	"github.com/nazzarr03/recipeApp/middleware"
	"github.com/nazzarr03/recipeApp/routes"
)

func main() {
	app := fiber.New()

	app.Use(middleware.LogMiddleware())

	app.Post("/signup", controllers.SignUp)
	app.Post("/login", controllers.Login)

	app.Use(middleware.Authentication())

	routes.RecipeRoutes(app)

	app.Listen(":3002")
}
