package responses

import (
	"prestamosbackend/initializers"
	"prestamosbackend/models"
)

type RouteResponse struct {
	ID      string          `json:"id" form:"id"`
	Name    string          `json:"name" form:"name"`
	Company CompanyResponse `json:"company" form:"company"`
}

func NewRouteResponse(route models.Route) *RouteResponse {

	db := initializers.DB

	var Company models.Company
	if err := db.First(&Company, route.CompanyID).Error; err != nil {
		return nil
	}

	return &RouteResponse{
		ID:      route.ID.String(),
		Name:    route.Name,
		Company: *NewCompanyResponse(Company),
	}
}

func NewRouteResponse1(route models.Route) RouteResponse {

	db := initializers.DB

	var Company models.Company
	if err := db.First(&Company, route.CompanyID).Error; err != nil {
		return RouteResponse{}
	}

	return RouteResponse{
		ID:      route.ID.String(),
		Name:    route.Name,
		Company: *NewCompanyResponse(Company),
	}
}
