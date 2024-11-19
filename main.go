package main

import (
	"log"
	"os"

	"finpocket.com/api/database"
	"finpocket.com/api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	database.ConnectDb()

	app := fiber.New()

	if err := routes.SeedCategories(); err != nil {
		log.Fatal("Error seeding categories:")
	}

	routes.Setup(app)

	log.Fatal(app.Listen(":" + port))
}
