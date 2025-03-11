package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"prestamosbackend/initializers"
	"time"
)

func Dashboard(c *gin.Context) {

	id := c.Param("id")

	currentDate := time.Now().Format("2006-01-02")

	// Monto Cobrado (Dinero) Listo
	var amountCollected float32
	initializers.DB.Table("fees").
		Select("COALESCE(SUM(paid_amount), 0)").
		Where("DATE(expected_date) = ? AND loan_id IN (SELECT id FROM loans WHERE route_id = ?)", currentDate, id).
		Scan(&amountCollected)

	// Porcentaje Cobrado Listo
	var percentageCharged float32
	var totalCustomers int64
	initializers.DB.Table("fees").
		Where("DATE(expected_date) = ? AND loan_id IN (SELECT id FROM loans WHERE route_id = ?)", currentDate, id).
		Count(&totalCustomers)
	if totalCustomers > 0 {
		var totalPaidCustomers int64
		initializers.DB.Table("fees").
			Where("DATE(expected_date) = ? AND loan_id IN (SELECT id FROM loans WHERE route_id = ?) AND paid_amount IS NOT NULL", currentDate, id).
			Count(&totalPaidCustomers)
		percentageCharged = (float32(totalPaidCustomers) / float32(totalCustomers)) * 100
	}

	// Prestamos nuevos (Cantidad) Listo
	var newLoansCount int64
	initializers.DB.Table("loans").
		Where("DATE(created_at) = ? AND route_id = ?", currentDate, id).
		Count(&newLoansCount)

	// Prestamos nuevos (Dinero) Listo
	var newLoansAmount float32
	initializers.DB.Table("loans").
		Select("COALESCE(SUM(amount), 0)").
		Where("DATE(created_at) = ? AND route_id = ?", currentDate, id).
		Scan(&newLoansAmount)

	// Cobros faltantes (cantidad) Listo
	var missingPaymentsCount int64
	initializers.DB.Table("fees").
		Where("DATE(expected_date) = ? AND loan_id IN (SELECT id FROM loans WHERE route_id = ?) AND paid_date IS NULL", currentDate, id).
		Count(&missingPaymentsCount)

	// Cobro Faltante (Dinero) con JOIN Listo
	var missingPaymentsAmount float32
	initializers.DB.Table("fees").
		Select("COALESCE(SUM((loans.interest / 100 * loans.amount) - paid_amount), 0)").
		Joins("INNER JOIN loans ON fees.loan_id = loans.id").
		Where("DATE(fees.expected_date) = ? AND loans.route_id = ?", currentDate, id).
		Scan(&missingPaymentsAmount)

	// Hora del primer Cobro Listo
	var firstPaymentTime *time.Time = nil
	if err := initializers.DB.Table("fees").
		Select("MIN(updated_at), NULL").
		Where("DATE(expected_date) = ? AND loan_id IN (SELECT id FROM loans WHERE route_id = ?)", currentDate, id).
		Scan(&firstPaymentTime); err != nil {
		firstPaymentTime = nil
	}

	// Hora del último Cobro Listo
	var lastPaymentTime *time.Time = nil
	if err := initializers.DB.Table("fees").
		Select("MAX(updated_at), NULL").
		Where("DATE(expected_date) = ? AND loan_id IN (SELECT id FROM loans WHERE route_id = ?)", currentDate, id).
		Scan(&lastPaymentTime); err != nil {
		lastPaymentTime = nil
	}

	type DashboardResponse struct {
		AmountCollected       string     `json:"amountCollected"`
		PercentageCharged     string     `json:"percentageCharged"`
		NewloansAmount        string     `json:"newloansAmount"`
		NewloansMoney         string     `json:"newloansMoney"`
		MissingPaymentsAmount string     `json:"missingPaymentsAmount"`
		MissingPaymentsMoney  string     `json:"missingPaymentsMoney"`
		FirstPaymentTime      *time.Time `json:"firstPaymentTime"`
		LastPaymentTime       *time.Time `json:"lastPaymentTime"`
	}

	amountCollectedStr := fmt.Sprintf("$ %.2f", amountCollected)
	percentageChargedStr := fmt.Sprintf("%.2f%%", percentageCharged)
	newLoansCountStr := fmt.Sprintf("%d", newLoansCount)
	newLoansAmountStr := fmt.Sprintf("$ %.2f", newLoansAmount)
	missingPaymentsCountStr := fmt.Sprintf("%d", missingPaymentsCount)
	missingPaymentsAmountStr := fmt.Sprintf("$ %.2f", missingPaymentsAmount)
	firstPaymentTimeStr := firstPaymentTime
	lastPaymentTimeStr := lastPaymentTime

	response := DashboardResponse{
		AmountCollected:       amountCollectedStr,
		PercentageCharged:     percentageChargedStr,
		NewloansAmount:        newLoansCountStr,
		NewloansMoney:         newLoansAmountStr,
		MissingPaymentsAmount: missingPaymentsCountStr,
		MissingPaymentsMoney:  missingPaymentsAmountStr,
		FirstPaymentTime:      firstPaymentTimeStr, // Obtén la hora del primer cobro
		LastPaymentTime:       lastPaymentTimeStr,  // Obtén la hora del último cobro
	}

	c.JSON(http.StatusOK, response)
}
