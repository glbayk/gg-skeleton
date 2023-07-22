package middleware

import (
	"strings"

	"github.com/glbayk/gg-skeleton/models"
	"github.com/glbayk/gg-skeleton/utils"
	"github.com/gofiber/fiber/v2"
)

func Authenticated(c *fiber.Ctx) error {
	var token string
	accessToken := c.Get("Authorization")
	if strings.HasPrefix(accessToken, "Bearer ") {
		token = strings.TrimPrefix(accessToken, "Bearer ")
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	claims, err := utils.TokenValidate(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	claimsEmail, err := claims.GetSubject()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	user := models.User{
		Email: claimsEmail,
	}
	err = user.Find()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	if user.ActivatedAt == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not Activated",
		})
	}

	c.Locals("user", user)

	return c.Next()
}
