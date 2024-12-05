package routes

import (
	"finpocket.com/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/users", CreateUser)
	app.Post("/users", CreateUser)

	auth := app.Group("", middleware.Auth)

	auth.Put("/users/picture", UpdateUserPicture)
	auth.Get("/plans", GetPlans)
	auth.Get("/plans/active", GetActivePlan)
	auth.Post("/plans", CreatePlan)
	auth.Put("/plans", UpdatePlan)
	auth.Delete("/plans", DisablePlan)

	auth.Get("/transactions", GetTransactions)
	auth.Post("/transactions", CreateTransaction)
	auth.Get("/transactions/summaries", GetTransactionSummaries)
}
