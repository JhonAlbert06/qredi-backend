package initializers

import (
	"fmt"
	"prestamosbackend/models"
)

func SyncDatabase() {
	err := DB.AutoMigrate(
		&models.SignatureType{},
		&models.Payment{},
		&models.Role{},
		&models.Company{},
		&models.Customer{},
		&models.Fee{},
		&models.Loan{},
		&models.Route{},
		&models.Spent{},
		&models.SpentType{},
		&models.CivilStatus{},
		&models.User{},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}
