package main

import (
	"log"
	"os"

	"finpocket.com/api/database"
	"finpocket.com/api/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	if err := routes.SeedCategories(); err != nil {
		log.Fatal("Error seeding categories:")
	}

	routes.Setup(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
