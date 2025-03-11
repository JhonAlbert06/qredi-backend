package responses

import (
	"prestamosbackend/initializers"
	"prestamosbackend/models"
)

type SpentResponse struct {
	ID     string             `json:"id" form:"id"`
	UserID string             `json:"userId" form:"userId"`
	Note   string             `json:"note" form:"note"`
	Cost   float32            `json:"cost" form:"cost"`
	Date   models.Date        `json:"date" form:"date"`
	Type   *SpentTypeResponse `json:"type" form:"type"`
}

func NewSpentResponse(spent models.Spent) *SpentResponse {

	db := initializers.DB

	var spentType models.SpentType
	if err := db.First(&spentType, spent.TypeId).Error; err != nil {
		return &SpentResponse{}
	}

	return &SpentResponse{
		ID:     spent.ID.String(),
		UserID: spent.UserId.String(),
		Note:   spent.Note,
		Cost:   spent.Cost,
		Date:   models.ToDate(spent.Date),
		Type:   NewSpentTypeResponse(spentType),
	}

}

func NewSpentResponse1(spent models.Spent) SpentResponse {

	db := initializers.DB

	var spentType models.SpentType
	if err := db.First(&spentType, spent.TypeId).Error; err != nil {
		return SpentResponse{}
	}

	return SpentResponse{
		ID:     spent.ID.String(),
		UserID: spent.UserId.String(),
		Note:   spent.Note,
		Cost:   spent.Cost,
		Date:   models.ToDate(spent.Date),
		Type:   NewSpentTypeResponse(spentType),
	}

}

type SpentTypeResponse struct {
	ID   string `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

func NewSpentTypeResponse(spentType models.SpentType) *SpentTypeResponse {
	return &SpentTypeResponse{
		ID:   spentType.ID.String(),
		Name: spentType.Name,
	}
}
