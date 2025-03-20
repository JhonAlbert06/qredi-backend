package controllers

import (
	"fmt"
	"net/http"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"prestamosbackend/responses"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateRoute(c *gin.Context) {
	routeID := uuid.New()

	var body struct {
		CompanyId string `json:"companyId" form:"companyId"`
		Name      string `json:"name" form:"name"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		fmt.Println(err)
		return
	}

	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Empty fields",
		})
		return
	}

	route := models.Route{
		ID:        routeID,
		Name:      body.Name,
		CompanyID: uuid.MustParse(body.CompanyId),
	}

	result := initializers.DB.Create(&route)

	if result.Error != nil {
		fmt.Println(result.Error)
		initializers.DB.Delete(&route)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create route",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func SearchRouteByParameter(c *gin.Context) {
	searchField := c.Query("field")

	user := c.MustGet("user").(models.User)

	var routes []models.Route

	db := initializers.DB

	db = db.Where("company_id = ?", user.CompanyId)

	db = db.Where("name LIKE ?", "%"+searchField+"%")

	if err := db.Find(&routes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error to load routes",
		})
		return
	}

	var routesResposne []responses.RouteResponse
	for _, route := range routes {
		routesResposne = append(routesResposne, responses.NewRouteResponse1(route))
	}

	if len(routesResposne) == 0 {
		c.JSON(http.StatusNotFound, 
			[]responses.RouteResponse{},
		)
		return
	}

	c.JSON(http.StatusOK, routesResposne)
}

func GetAllRoutes(c *gin.Context) {

	db := initializers.DB

	user := c.MustGet("user").(models.User)
	db = db.Where("company_id = ?", user.CompanyId)

	var routes []models.Route
	if err := db.Find(&routes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load route",
		})
		return
	}

	var routeResponse []responses.RouteResponse
	for _, route := range routes {
		routeResponse = append(routeResponse, responses.NewRouteResponse1(route))
	}

	c.JSON(http.StatusOK, routeResponse)

}

func SearchRouteById(c *gin.Context) {
	id := c.Param("id")

	var route models.Route

	if err := initializers.DB.Where("id = ?", id).Find(&route).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load route",
		})
		return
	}

	routeResponse := responses.NewRouteResponse(route)

	c.JSON(http.StatusOK, routeResponse)
}

func EditRoute(c *gin.Context) {
	var body struct {
		RouteID string `json:"id" form:"id"`
		Name    string `json:"name" form:"name"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		fmt.Println(err)
		return
	}

	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Empty fields",
		})
		return
	}

	var existingRoute models.Route

	if err := initializers.DB.Where("id = ?", body.RouteID).First(&existingRoute).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "Route not found",
		})
		return
	}

	user := c.MustGet("user").(models.User)

	if existingRoute.CompanyID != user.CompanyId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized",
		})
		return
	}

	existingRoute.Name = body.Name

	result := initializers.DB.Save(&existingRoute)

	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to update route",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
