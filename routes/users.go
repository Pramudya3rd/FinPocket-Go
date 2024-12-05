package routes

import (
	"finpocket.com/api/database"
	"finpocket.com/api/models"
	"finpocket.com/api/storage"
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

func UpdateUserPicture(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	file, err := c.FormFile("picture")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Jenis file harus berupa jpg atau png",
			"status":  "error",
			"data":    nil,
		})
	}

	blob, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}
	defer blob.Close()

	storage := storage.Init("users/pictures/")

	url, err := storage.UploadFile(blob, user.FirebaseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	user.Picture = url
	database.DBConn.Save(&user)

	return c.JSON(fiber.Map{
		"message": "User berhasil diupdate",
		"status":  "success",
		"data":    user,
	})
}
