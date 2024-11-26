package models

import "gorm.io/gorm"

// Model Transaksi
type Transaction struct {
	gorm.Model

	UserID      uint    `json:"user_id"`
	CategoryID  int64   `json:"category_id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Total       float64 `json:"total"`
}
