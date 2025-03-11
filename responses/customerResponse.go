package responses

import (
	"prestamosbackend/initializers"
	"prestamosbackend/models"
)

type CustomerResponse struct {
	ID          string          `json:"id" form:"id"`
	Company     CompanyResponse `json:"company" form:"company"`
	Cedula      string          `json:"cedula" form:"cedula"`
	FirstName   string          `json:"firstName" form:"firstName"`
	LastName    string          `json:"lastName" form:"lastName"`
	Address     string          `json:"address" form:"address"`
	Phone       string          `json:"phone" form:"phone"`
	CivilStatus string          `json:"civilStatus" form:"civilStatus"`
	Reference   string          `json:"reference" form:"reference"`
	Latitude    *string         `json:"latitude" form:"latitude"`
	Longitude   *string         `json:"longitude" form:"longitude"`
}

func NewCustomerResponse(customer models.Customer) *CustomerResponse {

	db := initializers.DB

	var company models.Company
	if err := db.First(&company, customer.CompanyId).Error; err != nil {
		return nil
	}

	return &CustomerResponse{
		ID:          customer.ID.String(),
		Company:     *NewCompanyResponse(company),
		Cedula:      customer.Cedula,
		FirstName:   customer.Names,
		LastName:    customer.LastNames,
		Address:     customer.Address,
		Phone:       customer.Phone,
		CivilStatus: customer.CivilStatus,
		Reference:   customer.Reference,
		Latitude:    customer.Latitude,
		Longitude:   customer.Longitude,
	}
}

func NewCustomerResponse1(customer models.Customer) CustomerResponse {

	db := initializers.DB

	var company models.Company
	if err := db.First(&company, customer.CompanyId).Error; err != nil {
		return CustomerResponse{}
	}

	return CustomerResponse{
		ID:          customer.ID.String(),
		Company:     *NewCompanyResponse(company),
		Cedula:      customer.Cedula,
		FirstName:   customer.Names,
		LastName:    customer.LastNames,
		Address:     customer.Address,
		Phone:       customer.Phone,
		CivilStatus: customer.CivilStatus,
		Reference:   customer.Reference,
		Latitude:    customer.Latitude,
		Longitude:   customer.Longitude,
	}
}
