package responses

import "prestamosbackend/models"

type CompanyResponse struct {
	ID   string `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

func NewCompanyResponse(company models.Company) *CompanyResponse {
	return &CompanyResponse{
		ID:   company.ID.String(),
		Name: company.Name,
	}
}
