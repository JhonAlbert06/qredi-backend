package responses

import (
	"github.com/google/uuid"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
)

type UserResponse struct {
	ID          uuid.UUID              `json:"id" form:"id"`
	CompanyID   uuid.UUID              `json:"companyId" form:"companyId"`
	IsAdmin     bool                   `json:"isAdmin" form:"isAdmin"`
	FirstName   string                 `json:"firstName" form:"firstName"`
	LastName    string                 `json:"lastName" form:"lastName"`
	UserName    string                 `json:"userName" form:"userName"`
	IsNew       bool                   `json:"isNew" form:"isNew"`
	Roles       []*RolesResponse       `json:"roles" form:"roles"`
	Memberships []*MembershipsResponse `json:"memberships" form:"memberships"`
}

func NewUserResponse(user models.User) *UserResponse {

	db := initializers.DB

	var company models.Company
	if err := db.Where("id = ?", user.CompanyId).First(&company).Error; err != nil {
		//company = nil
	}

	var UserRole models.UserRole
	if err := db.Where("user_id = ?", user.ID).First(&UserRole).Error; err != nil {
		//UserRole = nil
	}

	var RolesAccess []models.RoleAccess
	if err := db.Where("id = ?", UserRole.RoleAccessID).Find(&RolesAccess).Error; err != nil {

	}

	var Membership []models.Membership
	if err := db.Where("user_id = ?", user.ID).Find(&Membership).Error; err != nil {

	}

	var Roles []*RolesResponse
	for _, role := range RolesAccess {
		Roles = append(Roles, NewRolesResponse(role))
	}

	var Memberships []*MembershipsResponse
	for _, membership := range Membership {
		Memberships = append(Memberships, NewMembershipsResponse(membership))
	}

	var userResponse UserResponse
	userResponse.ID = user.ID
	userResponse.CompanyID = user.CompanyId
	userResponse.IsAdmin = user.IsAdmin
	userResponse.FirstName = user.FirstName
	userResponse.LastName = user.LastName
	userResponse.UserName = user.UserName
	userResponse.IsNew = user.IsNew
	userResponse.Roles = Roles
	userResponse.Memberships = Memberships

	return &userResponse
}
