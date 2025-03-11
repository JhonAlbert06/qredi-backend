package responses

import "prestamosbackend/models"

type GroupResponse struct {
	ID   string `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

func NewGroupResponse(group models.GroupRole) *GroupResponse {
	return &GroupResponse{
		ID:   group.ID.String(),
		Name: group.Name,
	}
}
