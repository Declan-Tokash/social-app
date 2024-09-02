package model
import (
 "github.com/google/uuid"
 "gorm.io/gorm"
)
// User struct
type User struct {
 gorm.Model
 ID       uuid.UUID `gorm:"type:uuid;"`
 Username string    `json:"username"`
 Email    string    `json:"email"`
 Password string    `json:"password"`
}
// Users struct
type Users struct {
 Users []User `json:"users"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
 // UUID version 4
 user.ID = uuid.New()
 return
}