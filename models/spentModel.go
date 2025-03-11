package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Spent struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key;"`
	CompanyID uuid.UUID `gorm:"type:char(36)"`
	UserId    uuid.UUID `gorm:"type:char(36)"`
	TypeId    uuid.UUID `gorm:"type:char(36)"`
	Note      string
	Cost      float32
	Date      time.Time
	gorm.Model
}
