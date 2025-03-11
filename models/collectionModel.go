package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Collection struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:char(36);primary_key;"`
	RouteId      uuid.UUID `gorm:"type:char(36)"`
	RouteName    string
	UserId       uuid.UUID `gorm:"type:char(36)"`
	UserName     string
	UserFullName string
	Date         time.Time
	Amount       float32
}
