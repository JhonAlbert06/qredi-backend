package responses

import (
	"prestamosbackend/initializers"
	"prestamosbackend/models"
)

type RolesResponse struct {
	ID    string        `json:"id" form:"id"`
	Name  string        `json:"name" form:"name"`
	Group GroupResponse `json:"group" form:"group"`
}

func NewRolesResponse(role models.RoleAccess) *RolesResponse {

	db := initializers.DB

	var group models.GroupRole
	if err := db.Where("id = ?", role.GroupID).First(&group).Error; err != nil {
		return nil
	}

	return &RolesResponse{
		ID:   role.ID.String(),
		Name: role.Name,
		Group: GroupResponse{
			ID:   group.ID.String(),
			Name: group.Name,
		},
	}
}
