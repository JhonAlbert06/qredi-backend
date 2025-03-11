package models

import "github.com/google/uuid"

type SignatureType struct {
	ID     uuid.UUID `gorm:"type:char(36);primary_key;"`
	LoanId uuid.UUID `gorm:"type:char(36)"`
	FeeId  uuid.UUID `gorm:"type:char(36)"`
	Name   string    `gorm:"unique"`
}
