package controllers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/recipe-app/config"
	"github.com/nazzarr03/recipe-app/models"
)

func GetRecipes(c *fiber.Ctx) error {
	var recipes []models.Recipe
	config.Db.Find(&recipes)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": recipes,
	})
}

func CreateRecipe(c *fiber.Ctx) error {
	var recipe models.Recipe
	if err := c.BodyParser(&recipe); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	recipe.CreatedAt = time.Now()

	if err := config.Db.Create(&recipe).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"data": recipe,
	})
}

func UpdateRecipe(c *fiber.Ctx) error {
	var incomingRecipe models.Recipe
	var recipe models.Recipe
	if err := c.BodyParser(&recipe); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	config.Db.First(&recipe, c.Params("id"))
	if recipe.ID == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Recipe not found",
		})
	}

	recipe.Title = incomingRecipe.Title
	recipe.Content = incomingRecipe.Content
	recipe.UpdatedAt = time.Now()

	if err := config.Db.Save(&recipe).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": recipe,
	})
}

func DeleteRecipe(c *fiber.Ctx) error {
	var recipe models.Recipe

	config.Db.First(&recipe, c.Params("id"))
	if recipe.ID == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Recipe not found",
		})
	}

	recipe.DeletedAt = time.Now()

	if err := config.Db.Save(&recipe).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Recipe deleted successfully",
	})
}
