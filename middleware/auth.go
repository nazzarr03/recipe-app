package middleware

import (
	"strings"

	"github.com/gofiber/fiber"
)

func Authentication(secret string) fiber.Handler {
	return func(c *fiber.Ctx) {
		authHeader := c.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) != 2 {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Not authorized",
			})
			c.Next()
			return
		}
		authToken := t[1]
		authorized, _ := IsAuthorized(authToken, secret)

		if !authorized {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Not authorized",
			})
			c.Next()
			return
		}

		ID, err := ExtractIDFromToken(authToken, secret)

		if err != nil {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Cannot extract user id from token",
			})
			c.Next()
			return
		}

		c.Set("id", ID)
		c.Next()
	}
}
