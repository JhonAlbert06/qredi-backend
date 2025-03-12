package responses

import (
	"prestamosbackend/models"
)

type RolesResponse struct {
	ID    string        `json:"id" form:"id"`
	Name  string        `json:"name" form:"name"`
}

func NewRolesResponse(role models.Role) *RolesResponse {

	//db := initializers.DB

	return &RolesResponse{
		ID:   role.ID.String(),
		Name: role.Name,
	}
}
