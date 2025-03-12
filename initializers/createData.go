package initializers

import (
	"prestamosbackend/models"

	"github.com/google/uuid"
)

func CreateData() {

	civilStatus := []models.CivilStatus{
		{ID: uuid.New(), Name: "Soltero"},
		{ID: uuid.New(), Name: "Casado"},
		{ID: uuid.New(), Name: "Divorciado"},
		{ID: uuid.New(), Name: "Viudo"},
	}

	for _, status := range civilStatus {
		DB.Create(&status)
	}

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
