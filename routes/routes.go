package routes

import (
	"finpocket.com/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/users", CreateUser)
	app.Put("/users/:user/picture", UpdateUserPicture)

	app.Get("/categories", GetCategories)

	auth := app.Group("", middleware.Auth)

	auth.Get("/plans", GetPlans)
	auth.Get("/plans/active", GetActivePlan)
	auth.Post("/plans", CreatePlan)
	auth.Put("/plans", UpdatePlan)
	auth.Delete("/plans", DisablePlan)

	auth.Get("/transactions", GetTransactions)
	auth.Post("/transactions", CreateTransaction)
	auth.Get("/transactions/summaries", GetTransactionSummaries)
}
