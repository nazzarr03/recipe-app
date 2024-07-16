package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/recipeApp/config"
	"github.com/nazzarr03/recipeApp/models"
	"github.com/streadway/amqp"
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

	recipe.Views = 0
	recipe.CreatedAt = time.Now()

	if err := config.Db.Create(&recipe).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var users []models.User
	config.Db.Find(&users)

	ch, err := config.RabbitMQConn.Channel()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
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
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	for _, user := range users {
		email := user.Email
		subject := "New Recipe Created!"
		body := fmt.Sprintf("Hello %s,\n\nA new recipe '%s' has been created. Check it out!\n\nBest regards,\nRecipe App Team", user.Email, recipe.Title)

		message := fmt.Sprintf("%s\n%s\n\n%s", email, subject, body)
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			},
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Failed to queue email for %s: %v", email, err),
			})
		}
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"data": recipe,
	})
}

func UpdateRecipe(c *fiber.Ctx) error {
	var incomingRecipe models.Recipe
	var recipe models.Recipe
	if err := c.BodyParser(&incomingRecipe); err != nil {
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

	if err := config.Db.Delete(&recipe).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := config.Rdb.Del(context.Background(), strconv.Itoa(int(recipe.ID))).Err(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Recipe deleted successfully",
	})
}

func GetRecipeByID(c *fiber.Ctx) error {
	var recipe models.Recipe

	config.Db.First(&recipe, c.Params("id"))
	if recipe.ID == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Recipe not found",
		})
	}

	recipe.Views++
	config.Db.Save(&recipe)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": recipe,
	})
}

func GetPopularRecipes(c *fiber.Ctx) error {
	var recipes []models.Recipe
	var popularRecipe []models.Recipe

	config.Db.Where("views > ?", 10).Find(&recipes)

	for _, recipe := range recipes {
		recipeJSON, err := json.Marshal(recipe)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		if err := config.Rdb.Set(context.Background(), strconv.Itoa(int(recipe.ID)), recipeJSON, 0).Err(); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	keys, err := config.Rdb.Keys(context.Background(), "*").Result()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	for _, key := range keys {
		value, err := config.Rdb.Get(context.Background(), key).Result()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		var recipe models.Recipe
		if err := json.Unmarshal([]byte(value), &recipe); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		popularRecipe = append(popularRecipe, recipe)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": popularRecipe,
	})
}
