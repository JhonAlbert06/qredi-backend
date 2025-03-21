package controllers

import (
	"fmt"
	"net/http"
	"prestamosbackend/initializers"
	"prestamosbackend/models"

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
	feeInterestAmount := ((float32(loan.Interest) / 100) * loan.Amount) + (loan.Amount / float32(loan.FeesQuantity))

	var payments []models.Payment
	if err := db.Where("fee_id = ?", fee.ID).Find(&payments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get payments",
		})
		return
	}

	var totalPaid float32
	for _, payment := range payments {
		totalPaid += payment.PaidAmount
	}

	if totalPaid+body.Amount > feeInterestAmount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The amount cannot be greater, the total amount to pay is " + fmt.Sprintf("%.2f", feeInterestAmount - totalPaid),
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

	fee.UpdatedAt = time.Now()
	if err := db.Save(&fee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update fee",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
