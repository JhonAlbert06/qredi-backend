package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	ID   uuid.UUID `gorm:"type:char(36);primary_key;"`
	Name string
	gorm.Model
}
