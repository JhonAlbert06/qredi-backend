package responses

import (
	"prestamosbackend/initializers"
	"prestamosbackend/models"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID          uuid.UUID              `json:"id" form:"id"`
	Company 		CompanyResponse				`json:"company" form:"company"`
	Roles 	RolesResponse					`json:"role" form:"role"`	
	FirstName   string                 `json:"firstName" form:"firstName"`
	LastName    string                 `json:"lastName" form:"lastName"`
	UserName    string                 `json:"userName" form:"userName"`
	IsNew       bool                   `json:"isNew" form:"isNew"`
}

func NewUserResponse(user models.User) *UserResponse {

	db := initializers.DB

	var company models.Company
	if err := db.Where("id = ?", user.CompanyId).First(&company).Error; err != nil {
		//company = nil
	}

	var role models.Role
	if err := db.Where("id = ?", user.RoleId).First(&role).Error; err != nil {
		//role = nil
	}


	var userResponse UserResponse
	userResponse.ID = user.ID
	userResponse.Company = *NewCompanyResponse(company)
	userResponse.Roles = *NewRolesResponse(role)
	userResponse.FirstName = user.FirstName
	userResponse.LastName = user.LastName
	userResponse.UserName = user.UserName
	userResponse.IsNew = user.IsNew

	return &userResponse
}
