package controllers

import (
	"errors"
	"net/http"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"prestamosbackend/responses"
	"prestamosbackend/utils"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


func CreateLoan(c *gin.Context) {
	var body struct {
		CustomerId   string  `form:"customerId" json:"customerId"`
		RouteId      string  `form:"routeId" json:"routeId"`
		Amount       float32 `form:"amount" json:"amount"`
		Interest     float32 `form:"interest" json:"interest"`
		FeesQuantity int32   `form:"feesQuantity" json:"feesQuantity"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Validar datos de entrada
	if body.Amount <= 0 || body.Interest < 0 || body.FeesQuantity < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan details"})
		return
	}

	if utils.HaveAnActiveLoan(body.CustomerId) {
		c.JSON(http.StatusConflict, gin.H{"message": "El cliente tiene un préstamo activo."})
		return
	}

	// Obtener usuario autenticado
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	user, ok := u.(models.User)
	if !ok || user.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	// Obtener cliente y ruta
	var customer models.Customer
	if err := initializers.DB.Where("id = ?", body.CustomerId).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	var route models.Route
	if err := initializers.DB.Where("id = ?", body.RouteId).First(&route).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	// Crear préstamo con transacción
	tx := initializers.DB.Begin()
	loanUUID := uuid.New()
	now := time.Now()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	loan := models.Loan{
		ID:            loanUUID,
		CustomerId:    customer.ID,
		RouteId:       route.ID,
		UserId:        user.ID,
		Amount:        body.Amount,
		Interest:      body.Interest,
		FeesQuantity:  body.FeesQuantity,
		Date:          currentDate,
		LoanIsPaid:    false,
		IsRenewed:     false,
		IsCurrentLoan: true,
	}

	if err := tx.Create(&loan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create loan"})
		return
	}

	// Crear cuotas (fees)
	for i := int32(1); i <= body.FeesQuantity; i++ {
		fee := models.Fee{
			ID:           uuid.New(),
			LoanId:       loanUUID,
			Number:       i,
			ExpectedDate: currentDate.Add(time.Duration(i) * 7 * 24 * time.Hour),
		}
		if err := tx.Create(&fee).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create loan fees"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, responses.NewLoanResponse(loan))
}

func SearchLoanById(c *gin.Context) {
	id := c.Param("id")

	var loan models.Loan
	if err := initializers.DB.Where("id = ?", id).First(&loan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to load loan",
			})
		}

		return
	}

	loanResponse := responses.NewLoanResponse(loan)

	c.JSON(http.StatusOK, loanResponse)
}

func CreateRenewLoan(c *gin.Context) {
	var body struct {
		LoanId       string  `form:"loanId" json:"loanId"`
		Amount       float32 `form:"amount" json:"amount"`
		Interest     float32 `form:"interest" json:"interest"`
		FeesQuantity int32   `form:"feesQuantity" json:"feesQuantity"`
	}

	db := initializers.DB

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Obtener el préstamo original
	var oldloan models.Loan
	if err := db.Where("id = ?", body.LoanId).First(&oldloan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	oldloan.IsCurrentLoan = false
	oldloan.LoanIsPaid = true
	if err := db.Save(&oldloan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update loan"})
		return
	}

	// Obtener cuotas del préstamo original
	var oldFees []models.Fee
	if err := db.Where("loan_id = ?", oldloan.ID).Find(&oldFees).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to load fees"})
		return
	}

	baseFee := oldloan.Amount / float32(oldloan.FeesQuantity) // Cuota base sin interés
  feeInterest := baseFee * (oldloan.Interest / 100)         // Interés de cada cuota
  oldfullAmount := (baseFee + feeInterest) * float32(oldloan.FeesQuantity)                       // Total por cuota con interés
	
	for _, fee := range oldFees {
		
		var payments []models.Payment
		if err := db.Where("fee_id = ?", fee.ID).Find(&payments).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get payments"})
			return
		}

		var totalPaid float32
		for _, payment := range payments {
			totalPaid += payment.PaidAmount
		}

		oldfullAmount -= totalPaid

	}


	// Crear nuevo préstamo
	newLoan := models.Loan{
		ID:            uuid.New(),
		CustomerId:    oldloan.CustomerId,
		RouteId:       oldloan.RouteId,
		UserId:        oldloan.UserId,
		Amount:        body.Amount + oldfullAmount,
		Interest:      body.Interest,
		FeesQuantity:  body.FeesQuantity,
		Date:          time.Now(),
		LoanIsPaid:    false,
		IsRenewed:     true,
		IsCurrentLoan: true,
	}

	if err := db.Create(&newLoan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create new loan"})
		return
	}

	// Crear cuotas (fees)
	for i := int32(1); i <= body.FeesQuantity; i++ {
		fee := models.Fee{
			ID:           uuid.New(),
			LoanId:       newLoan.ID,
			Number:       i,
			ExpectedDate: time.Now().Add(time.Duration(i) * 7 * 24 * time.Hour),
		}
		if err := db.Create(&fee).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create loan fees"})
			return
		}
	}

	c.JSON(http.StatusCreated, responses.NewLoanResponse(newLoan))
}
