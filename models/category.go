package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model

	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
