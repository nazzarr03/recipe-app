package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/recipe-app/controllers"
)

func RecipeRoutes(router *fiber.App) {
	router.Get("/recipes", controllers.GetRecipes)
	router.Post("/recipes", controllers.CreateRecipe)
	router.Put("/recipes/:id", controllers.UpdateRecipe)
	router.Delete("/recipes/:id", controllers.DeleteRecipe)
}
