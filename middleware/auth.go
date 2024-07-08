package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing JWT",
			})
		}

		token := authHeader[len("Bearer "):]

		claims, err := ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired JWT",
			})
		}

		c.Locals("userID", claims.UserId)

		return c.Next()
	}
}
