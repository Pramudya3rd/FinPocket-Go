package routes

import (
	"fmt"
	"math"

	"finpocket.com/api/database"
	"finpocket.com/api/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetPlans(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	plans := []models.Plan{}
	if err := database.DBConn.Where("user_id = ?", user.ID).Find(&plans).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	return c.JSON(plans)
}

func GetActivePlan(c *fiber.Ctx) error {
	month := c.Query("month")
	user := c.Locals("user").(*models.User)

	activePlan := models.Plan{}

	if month != "" {
		if err := database.DBConn.Where("user_id = ? AND MONTH(created_at) = ?", user.ID, month).
			Order("created_at DESC").
			First(&activePlan).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Active plan not found",
				"status":  "error",
				"data":    nil,
			})
		}
	} else {
		if err := database.DBConn.Where("user_id = ? AND active = ?", user.ID, true).
			First(&activePlan).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Active plan not found",
				"status":  "error",
				"data":    nil,
			})
		}
	}

	budgets := []models.Budget{}
	if err := database.DBConn.Where("plan_id = ? AND deleted_at IS NULL", activePlan.ID).
		Find(&budgets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	for idx, budget := range budgets {
		latestTransaction := models.Transaction{}
		database.DBConn.Where("category_id = ? AND MONTH(created_at) = ?", budget.CategoryID, activePlan.CreatedAt.Month()).
			Order("created_at DESC").
			First(&latestTransaction)
		if latestTransaction.ID != 0 {
			budgets[idx].Allocated = math.Abs(latestTransaction.Total)
		}
	}

	return c.JSON(fiber.Map{
		"plan":    activePlan,
		"budgets": budgets,
	})
}

type PlanRequest struct {
	Plan    models.Plan     `json:"plan" validate:"required"`
	Budgets []models.Budget `json:"budgets" validate:"required,min=1"`
}

func CreatePlan(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	validate := validator.New()

	p := new(PlanRequest)
	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	p.Plan.UserID = user.ID

	if err := validate.Struct(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	if p.Plan.Income <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Income must be greater than 0",
			"status":  "error",
			"data":    nil,
		})
	}

	for _, budget := range p.Budgets {
		if err := database.DBConn.First(&models.Category{}, budget.CategoryID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Category %d not found", budget.CategoryID),
				"status":  "error",
				"data":    nil,
			})
		} else if budget.Amount <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Budget amount must be greater than 0",
				"status":  "error",
				"data":    nil,
			})
		}
	}

	activePlan := models.Plan{}
	if err := database.DBConn.Where("user_id = ? AND active = ?", user.ID, true).
		First(&activePlan).Error; err == nil {
		activePlan.Active = false
		if err := database.DBConn.Save(&activePlan).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
				"status":  "error",
				"data":    nil,
			})
		}
	}

	if err := database.DBConn.Create(&p.Plan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	for i := range p.Budgets {
		p.Budgets[i].UserID = user.ID
		p.Budgets[i].PlanID = p.Plan.ID

		if err := database.DBConn.Create(&p.Budgets[i]).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
				"status":  "error",
				"data":    nil,
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Plan created successfully",
		"status":  "success",
		"data":    p,
	})
}

func UpdatePlan(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	validate := validator.New()

	p := new(PlanRequest)
	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	p.Plan.UserID = user.ID

	activePlan := models.Plan{}
	if err := database.DBConn.Where("user_id = ? AND active = ?", user.ID, true).
		First(&activePlan).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Active plan not found",
			"status":  "error",
			"data":    nil,
		})
	}

	if err := validate.Struct(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	if p.Plan.Income <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Income must be greater than 0",
			"status":  "error",
			"data":    nil,
		})
	}

	for _, budget := range p.Budgets {
		if err := database.DBConn.First(&models.Category{}, budget.CategoryID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Category %d not found", budget.CategoryID),
				"status":  "error",
				"data":    nil,
			})
		} else if budget.Amount <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Budget amount must be greater than 0",
				"status":  "error",
				"data":    nil,
			})
		}
	}

	activePlan.Income = p.Plan.Income
	activePlan.Dependents = p.Plan.Dependents

	if err := database.DBConn.Save(&activePlan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	if err := database.DBConn.Where("plan_id = ?", activePlan.ID).
		Delete(&models.Budget{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	for i := range p.Budgets {
		p.Budgets[i].UserID = user.ID
		p.Budgets[i].PlanID = activePlan.ID

		if err := database.DBConn.Create(&p.Budgets[i]).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
				"status":  "error",
				"data":    nil,
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Plan created successfully",
		"status":  "success",
		"data":    p,
	})
}

func DisablePlan(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	activePlan := models.Plan{}
	if err := database.DBConn.Where("user_id = ? AND active = ?", user.ID, true).
		First(&activePlan).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Active plan not found",
			"status":  "error",
			"data":    nil,
		})
	}

	activePlan.Active = false
	if err := database.DBConn.Save(&activePlan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	return c.JSON(user)
}
