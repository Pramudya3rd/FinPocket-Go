package routes

import "github.com/gofiber/fiber/v2"

func Setup(app *fiber.App) {
	app.Post("/users", CreateUser)
	app.Put("/users/:user/picture", UpdateUserPicture)

	app.Get("/categories", GetCategories)
}
