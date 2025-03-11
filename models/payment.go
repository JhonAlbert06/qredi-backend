package models

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID         uuid.UUID `gorm:"type:char(36);primary_key;"`
	FeeId      uuid.UUID `gorm:"type:char(36)"`
	PaidAmount float32
	PaidDate   time.Time
	UserId     uuid.UUID `gorm:"type:char(36)"`
}
