package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleAccess struct {
	ID      uuid.UUID `gorm:"type:char(36);primary_key;"`
	Name    string
	GroupID uuid.UUID `gorm:"type:char(36)"`
	gorm.Model
}
