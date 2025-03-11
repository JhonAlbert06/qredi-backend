package initializers

import (
	"github.com/google/uuid"
	"prestamosbackend/models"
)

func CreateData() {
	spentType := []models.SpentType{
		{ID: uuid.New(), Name: "Combustible"},
		{ID: uuid.New(), Name: "Dieta"},
		{ID: uuid.New(), Name: "Refraccion"},
		{ID: uuid.New(), Name: "Otro"},
	}

	for _, Stype := range spentType {
		DB.Create(&Stype)
	}
}
