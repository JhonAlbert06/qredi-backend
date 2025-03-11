package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole struct {
	ID           uuid.UUID `gorm:"type:char(36);primary_key;"`
	RoleAccessID uuid.UUID `gorm:"type:char(36)"`
	UserID       uuid.UUID `gorm:"type:char(36)"`
	gorm.Model
}
