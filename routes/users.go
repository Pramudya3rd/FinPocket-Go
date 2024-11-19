package routes

import (
	"finpocket.com/api/database"
	"finpocket.com/api/models"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	exists := database.DBConn.Where("firebase_id = ?", user.FirebaseID).FirstOrCreate(&user)
	if err := exists.Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"message": "User berhasil ditambahkan",
		"status":  "success",
		"data":    user,
	})
}
