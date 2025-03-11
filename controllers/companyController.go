package controllers

import (
	"net/http"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"prestamosbackend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCompany(c *gin.Context) {

	var body struct {
		Name string `json:"name" form:"name"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// IsStringEmpty checks if a string is null or empty.
	if utils.IsStringEmpty(body.Name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name cannot be empty",
		})
		return
	}

	// Create the company
	company := models.Company{
		ID:   uuid.New(),
		Name: body.Name,
	}

	// Save the company
	db := initializers.DB
	if err := db.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create company",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"company": company,
	})
}
