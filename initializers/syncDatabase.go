package initializers

import (
	"fmt"
	"prestamosbackend/models"
)

func SyncDatabase() {
	err := DB.AutoMigrate(
		&models.SignatureType{},
		&models.Payment{},
		&models.UserRole{},
		&models.RoleAccess{},
		&models.Membership{},
		&models.Company{},
		&models.Customer{},
		&models.Fee{},
		&models.Loan{},
		&models.Route{},
		&models.Spent{},
		&models.SpentType{},
		&models.User{},
		&models.Collection{},
		&models.CollectionDetail{},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}
