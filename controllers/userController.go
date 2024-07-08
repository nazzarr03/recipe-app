package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/recipe-app/config"
	"github.com/nazzarr03/recipe-app/middleware"
	"github.com/nazzarr03/recipe-app/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if user.Username == "" || user.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Username and password are required",
		})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	if err := config.Db.Create(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user,
	})

}

func Login(c *fiber.Ctx) error {
	var user models.User
	var existingUser models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if user.Username == "" || user.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Username and password are required",
		})
	}

	config.Db.Where("username = ?", user.Username).First(&existingUser)
	if existingUser.ID == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	secret := os.Getenv("JWT_SECRET")
	accessToken, err := middleware.CreateAccessToken(&user, secret, 24)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":      "Login successful",
		"access_token": accessToken,
	})
}
