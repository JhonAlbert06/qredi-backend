package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"prestamosbackend/responses"
)

func CreateTypeSpent(c *gin.Context) {
	var body struct {
		Name string `json:"name"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		fmt.Println(err)
		return
	}

	typeId := uuid.New()
	typeS := models.SpentType{
		ID:   typeId,
		Name: body.Name,
	}

	var db = initializers.DB

	result := db.Create(&typeS)

	if result.Error != nil {
		fmt.Println(result.Error)
		db.Unscoped().Delete(&typeS)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create spent type",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func GetAllTypesSpent(c *gin.Context) {
	var spentType []models.SpentType
	var db = initializers.DB

	result := db.Find(&spentType)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve spents",
		})
		return
	}

	var typeResponse []*responses.SpentTypeResponse
	for _, Type := range spentType {
		typeResponse = append(typeResponse, responses.NewSpentTypeResponse(Type))
	}

	c.JSON(http.StatusOK, typeResponse)
}
