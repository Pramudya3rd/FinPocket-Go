package routes

import (
	"finpocket.com/api/database"
	"finpocket.com/api/models"
	"github.com/gofiber/fiber/v2"
)

// Seed data untuk kategori
var Categories = []models.Category{
	{Name: "Bills"},
	{Name: "Groceries"},
	{Name: "Transport"},
	{Name: "Entertainments"},
	{Name: "Healthcare"},
	{Name: "Education"},
	{Name: "Utility"},
	{Name: "Saving"},
}

// Fungsi untuk seed data ke tabel categories
func SeedCategories() error {
	var count int64
	if err := database.DBConn.Model(&models.Category{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		if err := database.DBConn.Create(&Categories).Error; err != nil {
			return err
		}
	}
	return nil
}

// Handler untuk mendapatkan data kategori
func GetCategories(c *fiber.Ctx) error {
	var categories []models.Category
	if err := database.DBConn.Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data",
			"status":  "error",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data berhasil diambil",
		"status":  "success",
		"data":    categories,
	})
}
