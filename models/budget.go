package models

import "gorm.io/gorm"

type Budget struct {
	gorm.Model

	UserID     uint    `json:"user_id" gorm:"not null"`
	PlanID     uint    `json:"plan_id" gorm:"not null"`
	CategoryID uint    `json:"category_id" gorm:"not null"`
	Percentage float32 `json:"percentage" gorm:"not null"`
	Amount     float64 `json:"amount" gorm:"not null"`

	Allocated float64 `json:"allocated" gorm:"-"`
}
