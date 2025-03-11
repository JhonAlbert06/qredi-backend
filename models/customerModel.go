package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key;"`
	CompanyId   uuid.UUID `gorm:"type:char(36)"`
	Cedula      string    `gorm:"unique"`
	Names       string
	LastNames   string
	Address     string
	Phone       string
	CivilStatus string
	Reference   string
	Latitude    *string
	Longitude   *string
	gorm.Model
}
