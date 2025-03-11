package utils

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"strings"
)

func FormatearCedula(input string) string {

	if len(input) < 3 {
		return input
	}

	parte1 := input[:3]
	parte2 := ""
	parte3 := ""

	if len(input) >= 10 {
		parte2 = input[3:10]
	} else {
		parte2 = input[3:]
	}

	if len(input) >= 11 {
		parte3 = input[10:]
	}

	resultado := parte1 + "-" + parte2

	if parte3 != "" {
		resultado += "-" + parte3
	}

	return resultado
}

func FormatearTelefono(input string) string {

	if len(input) < 3 {
		return input
	}

	parte1 := input[:3]
	parte2 := ""
	parte3 := ""

	if len(input) >= 6 {
		parte2 = input[3:6]
	}

	if len(input) >= 10 {
		parte3 = input[6:10]
	}

	resultado := parte1

	if parte2 != "" {
		resultado += "-" + parte2
	}

	if parte3 != "" {
		resultado += "-" + parte3
	}

	return resultado
}

func EliminarGuiones(input string) string {
	return strings.ReplaceAll(input, "-", "")
}

func HaveAnActiveLoan(customerId string) bool {
	db := initializers.DB
	var count int64
	err := db.Table("loans").Where("customer_id = ? AND loan_is_paid = ?", customerId, false).Count(&count).Error
	if err != nil {
		return false
	}

	return count > 0
}

func HaveanactiveloanRenew(customerId string) bool {
	db := initializers.DB
	var count int64
	err := db.Table("loans").Where("customer_id = ? AND loan_is_paid = ? AND is_current_loan = ?", customerId, true, true).Count(&count).Error
	if err != nil {
		return false
	}

	return count > 0
}

func MarkLoanAsPaidIfAllFeesPaid(db *gorm.DB, loanID uuid.UUID) bool {
	var loan models.Loan
	if err := db.First(&loan, loanID).Error; err != nil {
		return false
	}

	// Calcula el interés esperado basado en los parámetros del préstamo
	expectedInterest := (float32(loan.Interest) / 100) * loan.Amount * float32(loan.FeesQuantity)
	var totalPaid float32

	var fees []models.Fee
	if err := db.Where("loan_id = ?", loanID).Find(&fees).Error; err != nil {
		return false
	}

	// Calcula la suma de los montos pagados en todas las cuotas
	for _, fee := range fees {

		var payments []models.Payment
		if err := db.Where("fee_id = ?", fee.ID).Find(&payments).Error; err != nil {
			return false
		}

		for _, payment := range payments {
			totalPaid += payment.PaidAmount
		}

	}

	// Comprueba si el préstamo está pagado
	if totalPaid == expectedInterest {
		loan.LoanIsPaid = true
		if err := db.Save(&loan).Error; err != nil {
			return false
		}
		return true
	}

	return false
}
