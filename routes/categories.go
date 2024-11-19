package routes

import (
	"finpocket.com/api/database"
	"finpocket.com/api/models"
	"github.com/gofiber/fiber/v2"
)

var Categories = []models.Category{
	{ID: "1", Name: "Bills"},
	{ID: "2", Name: "Groceries"},
	{ID: "3", Name: "Transport"},
	{ID: "4", Name: "Entertainments"},
	{ID: "5", Name: "Healthcare"},
	{ID: "6", Name: "Education"},
	{ID: "7", Name: "Utility"},
	{ID: "8", Name: "Saving"},
}

func SeedCategories() error {
	var count int64
	database.DBConn.Model(&models.Category{}).Count(&count)

	if count == 0 {
		if err := database.DBConn.Create(&Categories).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetCategories(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Data berhasil diambil",
		"status":  "success",
		"data":    Categories,
	})
}
