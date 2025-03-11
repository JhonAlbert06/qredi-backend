package responses

import "prestamosbackend/models"

type MembershipsResponse struct {
	ID    string      `json:"id" form:"id"`
	Date  models.Date `json:"date" form:"date"`
	Days  int         `json:"days" form:"days"`
	Price float32     `json:"price" form:"price"`
}

func NewMembershipsResponse(membership models.Membership) *MembershipsResponse {
	return &MembershipsResponse{
		ID:    membership.ID.String(),
		Date:  models.ToDate(membership.Date),
		Price: membership.Price,
		Days:  membership.Days,
	}
}
