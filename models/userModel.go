package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID `gorm:"type:char(36);primary_key;"`
	CompanyId       uuid.UUID `gorm:"type:char(36)"`
	IsAdmin         bool
	IsNew           bool
	FirstName       string
	LastName        string
	UserName        string `gorm:"unique"`
	Password        string
	PasswordVersion uuid.UUID `gorm:"type:char(36)"`
	gorm.Model
}
