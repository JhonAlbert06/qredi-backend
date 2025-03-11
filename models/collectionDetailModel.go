package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectionDetail struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:char(36);primary_key;"`
	CollectionId     uuid.UUID `gorm:"type:char(36)"`
	LoanId           uuid.UUID `gorm:"type:char(36)"`
	FeeId            uuid.UUID `gorm:"type:char(36)"`
	FeeNumber        int32
	FeeQuantity      int32
	ExpectedAmount   float32
	PaidAmount       float32
	CustomerId       uuid.UUID `gorm:"type:char(36)"`
	CustomerFullName string
}
