package initializers

import (
	"prestamosbackend/models"

	"github.com/google/uuid"
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

	role := []models.Role{
		{ID: uuid.New(), Name: "Admin"},
		{ID: uuid.New(), Name: "User"},
	}

	for _, r := range role {
		DB.Create(&r)
	}

}
