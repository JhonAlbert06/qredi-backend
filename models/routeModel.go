package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Route struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key;"`
	CompanyID uuid.UUID `gorm:"type:char(36)"`
	Name      string    `gorm:"unique"`
	gorm.Model
}
