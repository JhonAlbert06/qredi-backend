package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"prestamosbackend/responses"
	"time"
)

func CreateSpent(c *gin.Context) {
	var body struct {
		CompanyId string  `json:"companyId" form:"companyId"`
		TypeId    string  `json:"typeId" form:"typeId"`
		Note      string  `json:"note" form:"note"`
		Cost      float32 `json:"cost" form:"cost"`
	}

	Date := time.Now()

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		fmt.Println(err)
		return
	}

	//Get the user
	u, _ := c.Get("user")
	if u.(models.User).ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user",
		})
		return
	}

	// Parse the user
	var user models.User
	if u, ok := u.(models.User); ok {
		user = u
	} else {
		fmt.Println(ok)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to load user",
		})
	}

	if body.CompanyId == "" || body.TypeId == "" || body.Note == "" || body.Cost == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid parameters",
		})
		return
	}

	spentId := uuid.New()
	spent := models.Spent{
		ID:        spentId,
		CompanyID: uuid.MustParse(body.CompanyId),
		UserId:    user.ID,
		TypeId:    uuid.MustParse(body.TypeId),
		Cost:      body.Cost,
		Note:      body.Note,
		Date:      Date,
	}

	var db = initializers.DB

	result := db.Create(&spent)

	if result.Error != nil {
		fmt.Println(result.Error)
		db.Unscoped().Delete(&spent)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create spent",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func GetSpent(c *gin.Context) {
	spentId := c.Param("id")
	var spent models.Spent

	var db = initializers.DB

	user := c.MustGet("user").(models.User)

	db.Where("company_id = ?", user.CompanyId)
	result := db.First(&spent).Where("id = ?", spentId)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Spent not found",
		})
		return
	}

	c.JSON(http.StatusOK, responses.NewSpentResponse(spent))
}

func GetAllSpents(c *gin.Context) {
	var spents []models.Spent
	var db = initializers.DB

	result := db.Order("created_at desc").Find(&spents)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve spents",
		})
		return
	}

	var spentsResponse []responses.SpentResponse
	for _, spent := range spents {

		var types models.SpentType
		result = db.First(&types, spent.TypeId)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Type not found",
			})
			return
		}

		var user models.User
		result = db.First(&user).Where("id = ?", spent.UserId)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		spentsResponse = append(spentsResponse, responses.NewSpentResponse1(spent))
	}

	c.JSON(http.StatusOK, spentsResponse)
}
