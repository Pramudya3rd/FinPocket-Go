package middleware

import (
	"strings"

	"finpocket.com/api/database"
	"finpocket.com/api/models"
	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	if c.Get("Authorization") == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"status:": "error",
		})
	}

	header := strings.Split(c.Get("Authorization"), " ")
	if len(header) != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"status:": "error",
		})
	} else if header[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"status:": "error",
		})
	}

	firebaseId := header[1]

	var user models.User
	database.DBConn.Where("firebase_id = ?", firebaseId).First(&user)

	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"status:": "error",
		})
	}

	c.Locals("user", &user)

	return c.Next()
}
