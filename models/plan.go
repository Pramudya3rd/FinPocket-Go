package models

import (
	"time"

	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model

	UserID     uint       `json:"user_id" validate:"required" gorm:"not null"`
	Income     float64    `json:"income" validate:"required" gorm:"not null"`
	Age        uint8      `json:"age" validate:"required,min=1" gorm:"not null"`
	Dependents uint8      `json:"dependents" gorm:"not null"`
	Active     bool       `json:"active" gorm:"default:true"`
	DisabledAt *time.Time `json:"disabled_at" gorm:"default:null"`
}
