package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Picture    string `json:"picture"`
	FirebaseID string `json:"firebase_id" gorm:"unique"`
}
