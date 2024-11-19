package main

import (
	"log"

	"finpocket.com/api/database"
	"finpocket.com/api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDb()
	app := fiber.New()

	if err := routes.SeedCategories(); err != nil {
		log.Fatal("Error seeding categories:")
	}

	routes.Setup(app)

	log.Fatal(app.Listen(":3000"))
}
