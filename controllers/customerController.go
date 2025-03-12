package controllers

import (
	"fmt"
	"net/http"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"prestamosbackend/responses"
	"prestamosbackend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCustomer(c *gin.Context) {

	// Create a new customer
	customerID := uuid.New()

	//Get the body
	var body struct {
		CompanyId   string  `form:"companyId" json:"companyId"`
		Cedula      string  `form:"cedula" json:"cedula"`
		Names       string  `form:"names" json:"names"`
		LastNames   string  `form:"lastNames" json:"lastNames"`
		Address     string  `form:"address" json:"address"`
		Phone       string  `form:"phone" json:"phone"`
		CivilStatusId string  `form:"civilStatusId" json:"civilStatusId"`
		Reference   string  `form:"reference" json:"reference"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		fmt.Println(err)
		return
	}

	if body.Cedula == "" || body.Names == "" || body.LastNames == "" || body.Address == "" || body.Phone == "" || body.CivilStatusId == "" || body.Reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Empty fields",
		})
		return
	}

	companyId := uuid.MustParse(body.CompanyId)

	// Create Customer with file paths
	customer := models.Customer{
		CompanyId:   companyId,
		ID:          customerID,
		Cedula:      body.Cedula,
		Names:       body.Names,
		LastNames:   body.LastNames,
		Address:     body.Address,
		Phone:       body.Phone,
		CivilStatusId: body.CivilStatusId,
		Reference:   body.Reference,
	}

	result := initializers.DB.Create(&customer)

	if result.Error != nil {
		fmt.Println(result.Error)
		initializers.DB.Unscoped().Delete(&customer)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create customer",
		})
		return
	}

	customerRes := responses.NewCustomerResponse(customer)

	// Respond
	c.JSON(http.StatusCreated, customerRes)
}

func UpdateCustomer(c *gin.Context) {

	var body struct {
		Id          string  `form:"id" json:"id"`
		Cedula      string  `form:"cedula" json:"cedula"`
		Names       string  `form:"names" json:"names"`
		LastNames   string  `form:"lastNames" json:"lastNames"`
		Address     string  `form:"address" json:"address"`
		Phone       string  `form:"phone" json:"phone"`
		CivilStatusId string  `form:"civilStatusId" json:"civilStatusId"`
		Reference   string  `form:"reference" json:"reference"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		fmt.Println(err)
		return
	}

	var existingCustomer models.Customer
	if err := initializers.DB.Where("id = ?", body.Id).First(&existingCustomer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Customer not found",
		})
		return
	}

	if body.Cedula != "" {
		existingCustomer.Cedula = body.Cedula
	}

	if body.Names != "" {
		existingCustomer.Names = body.Names
	}

	if body.LastNames != "" {
		existingCustomer.LastNames = body.LastNames
	}

	if body.Address != "" {
		existingCustomer.Address = body.Address
	}

	if body.Phone != "" {
		existingCustomer.Phone = body.Phone
	}

	if body.CivilStatusId != "" {
		existingCustomer.CivilStatusId = body.CivilStatusId
	}

	if body.Reference != "" {
		existingCustomer.Reference = body.Reference
	}
	
	result := initializers.DB.Save(&existingCustomer)
	if result.Error != nil {
		fmt.Println(result.Error)
		initializers.DB.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update customer",
		})
		return
	}

	// Respond customer
	c.JSON(http.StatusOK, responses.NewCustomerResponse(existingCustomer))
}

func SearchCustomerByParameter(c *gin.Context) {
	query := c.Query("query")
	searchField := c.Query("field")

	db := initializers.DB

	user := c.MustGet("user").(models.User)
	db = db.Where("company_id = ?", user.CompanyId)

	var customers []models.Customer

	switch searchField {
	case "names":
		db = db.Where("names LIKE ?", "%"+query+"%").Limit(15)
	case "last_names":
		db = db.Where("last_names LIKE ?", "%"+query+"%").Limit(15)
	case "cedula":
		query = utils.EliminarGuiones(query)
		query = utils.FormatearCedula(query)
		db = db.Where("cedula LIKE ?", "%"+query+"%").Limit(15)
	case "phone":
		query = utils.EliminarGuiones(query)
		query = utils.FormatearTelefono(query)
		db = db.Where("phone LIKE ?", "%"+query+"%").Limit(15)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Campo de búsqueda no válido",
		})
		return
	}

	if err := db.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error to load customers",
		})
		return
	}

	var customerResponse []responses.CustomerResponse

	for _, customer := range customers {
		customerRes := responses.NewCustomerResponse1(customer)
		customerResponse = append(customerResponse, customerRes)
	}

	c.JSON(http.StatusOK, customerResponse)
}

func SearchCustomerById(c *gin.Context) {
	id := c.Param("id")

	db := initializers.DB

	user := c.MustGet("user").(models.User)
	db = db.Where("company_id = ?", user.CompanyId)

	var customer models.Customer
	if err := db.Where("id = ?", id).Find(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load customer",
		})
		return
	}

	customerResponse := responses.NewCustomerResponse(customer)

	c.JSON(http.StatusOK, customerResponse)
}
