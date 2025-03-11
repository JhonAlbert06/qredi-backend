package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Fee struct {
	ID           uuid.UUID `gorm:"type:char(36);primary_key;"`
	Number       int32
	LoanId       uuid.UUID `gorm:"type:char(36)"`
	ExpectedDate time.Time
	gorm.Model
}
