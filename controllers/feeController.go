package controllers

import (
	"fmt"
	"net/http"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"prestamosbackend/responses"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PayOffFee(c *gin.Context) {
	var body struct {
		ID     string  `json:"id" form:"id"`
		Amount float32 `json:"amount" form:"amount"`
	}

	var db = initializers.DB

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		fmt.Println(err)
		return
	}

	var fee models.Fee
	if err := db.Where("id = ?", body.ID).First(&fee).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "cuota no encontrada",
		})
		return
	}

	var loan models.Loan
	if err := db.Where("id = ?", fee.LoanId).First(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "prestamo no encontrado",
		})
		return
	}

	user := c.MustGet("user").(models.User)
	interestAmount := float32(loan.Interest) / 100 * loan.Amount

	if body.Amount > float32(int(interestAmount)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The amount cannot be greater",
		})
		return
	}

	var payment = models.Payment{
		ID:         uuid.New(),
		UserId:     user.ID,
		FeeId:      fee.ID,
		PaidDate:   time.Now(),
		PaidAmount: body.Amount,
	}

	if err := db.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create payment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func GetFeesByDate(c *gin.Context) {

	date := c.Query("date")
	routeId := c.Query("routeId")

	db := initializers.DB

	var fees []models.Fee
	if err := db.Table("fees").
		Select("fees.*").
		Joins("JOIN loans ON loans.id = fees.loan_id").
		Where("fees.expected_date = ? AND loans.route_id = ?", date, routeId).
		Find(&fees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get fees",
		})
	}

	sort.Slice(fees, func(i, j int) bool {
		return fees[i].Number < fees[j].Number
	})

	var feesResponse []*responses.FeeResponse
	for _, fee := range fees {
		feesResponse = append(feesResponse, responses.NewFeeResponse(fee))
	}

	c.JSON(http.StatusOK, feesResponse)
}
