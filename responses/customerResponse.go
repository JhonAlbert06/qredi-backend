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
}

func NewCustomerResponse(customer models.Customer) *CustomerResponse {

	db := initializers.DB

	var company models.Company
	if err := db.First(&company, customer.CompanyId).Error; err != nil {
		return nil
	}

	var civilStatus models.CivilStatus
	if err := db.Where("id = ?", customer.CivilStatusId).First(&civilStatus).Error; err != nil {
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
		CivilStatus: civilStatus.Name,
		Reference:   customer.Reference,
	}
}

func NewCustomerResponse1(customer models.Customer) CustomerResponse {

	db := initializers.DB

	var company models.Company
	if err := db.First(&company, customer.CompanyId).Error; err != nil {
		return CustomerResponse{}
	}

	var civilStatus models.CivilStatus
	if err := db.First(&civilStatus, customer.CivilStatusId).Error; err != nil {
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
		CivilStatus: civilStatus.Name,
		Reference:   customer.Reference,
	}
}
