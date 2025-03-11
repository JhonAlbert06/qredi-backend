package models

import (
	"github.com/google/uuid"
	"time"
)

type Membership struct {
	ID     uuid.UUID `gorm:"type:char(36);primary_key;"`
	UserId uuid.UUID `gorm:"type:char(36)"`
	Date   time.Time
	Price  float32
	Days   int
}
