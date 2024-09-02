package model

import (
	"gorm.io/gorm"
   )

type Location struct {
	gorm.Model
	UserID    string  `json:userID`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}