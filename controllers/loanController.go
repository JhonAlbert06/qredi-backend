package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"prestamosbackend/responses"
	"prestamosbackend/utils"
	"sort"
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

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	if utils.HaveAnActiveLoan(body.CustomerId) {
		c.JSON(http.StatusConflict, gin.H{
			"message": "El cliente tiene un préstamo activo.",
		})
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

	var customer models.Customer
	if err := initializers.DB.Where("id = ?", body.CustomerId).Find(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load customer",
		})
		return
	}

	var route models.Route
	if err := initializers.DB.Where("id = ?", body.RouteId).Find(&route).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load route",
		})
		return
	}

	loanUUID := uuid.New()
	now := time.Now()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	loan := models.Loan{
		Model:         gorm.Model{},
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

	resultLoan := initializers.DB.Create(&loan)
	if resultLoan.Error != nil {
		fmt.Println(resultLoan.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create loan",
		})
		return
	}

	for i := int32(1); i <= body.FeesQuantity; i++ {
		fee := models.Fee{ID: uuid.New(), LoanId: loanUUID, Number: i, ExpectedDate: currentDate.Add(time.Duration(i) * 7 * 24 * time.Hour)}
		resultFee := initializers.DB.Create(&fee)
		if resultFee.Error != nil {
			fmt.Println(resultFee.Error)
			initializers.DB.Unscoped().Delete(&loan)
			initializers.DB.Unscoped().Where("loan_id = ?", loanUUID).Delete(models.Fee{})
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create loan",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, responses.NewLoanResponse(loan))
}

func SearchLoanByParameter(c *gin.Context) {

	query := c.Query("query")
	searchField := c.Query("field")

	switch searchField {
	case "names":
	case "last_names":
	case "cedula":
		query = utils.EliminarGuiones(query)
		query = utils.FormatearCedula(query)
	case "phone":
		query = utils.EliminarGuiones(query)
		query = utils.FormatearTelefono(query)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Campo de búsqueda no válido",
		})
		return
	}

	var db = initializers.DB

	user1 := c.MustGet("user").(models.User)

	var loans []models.Loan
	if err := db.Table("loans").
		Select("loans.*").
		Joins("JOIN customers ON loans.customer_id = customers.id AND customers.company_id = ?", user1.CompanyId).
		Where("customers."+searchField+" LIKE ?", "%"+query+"%").
		Where("loans.loan_is_paid = ?", false).
		Find(&loans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load loans",
		})
		return
	}

	if len(loans) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No found",
		})
		return
	}

	var loansResponse []responses.LoanResponse

	for _, loan := range loans {

		loansResponse = append(loansResponse, responses.NewLoanResponse1(loan))
	}

	c.JSON(http.StatusOK, loansResponse)
}

func GetLoansByDate(c *gin.Context) {

	date := c.Query("date")
	routeId := c.Query("routeId")

	var loans []models.Loan
	user1 := c.MustGet("user").(models.User)

	if err := initializers.DB.Table("loans").
		Select("DISTINCT loans.*").
		Joins("JOIN fees ON loans.id = fees.loan_id").
		Joins("JOIN customers ON loans.customer_id = customers.id AND customers.company_id = ?", user1.CompanyId).
		Where("fees.expected_date <= ?", date).
		Where("loans.loan_is_paid = 0").
		Where("loans.route_id = ?", routeId).
		Find(&loans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load loans",
		})
		return
	}

	if len(loans) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No found",
		})
		return
	}

	var loansResponse []responses.LoanResponse

	for _, loan := range loans {

		loansResponse = append(loansResponse, responses.NewLoanResponse1(loan))
	}

	c.JSON(http.StatusOK, loansResponse)
}

func GetLoansByParameter(c *gin.Context) {

	date := c.Query("date")
	routeId := c.Query("routeId")
	query := c.Query("query")
	searchField := c.Query("field")

	switch searchField {
	case "names":
	case "last_names":
	case "cedula":
		query = utils.EliminarGuiones(query)
		query = utils.FormatearCedula(query)
	case "phone":
		query = utils.EliminarGuiones(query)
		query = utils.FormatearTelefono(query)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Campo de búsqueda no válido",
		})
		return
	}

	var loans []models.Loan

	user1 := c.MustGet("user").(models.User)

	if err := initializers.DB.Table("loans").
		Select("DISTINCT loans.*").
		Joins("JOIN fees ON loans.id = fees.loan_id").
		Joins("JOIN customers ON loans.customer_id = customers.id AND customers.company_id = ?", user1.CompanyId).
		Where("fees.expected_date <= ?", date).
		Where("loans.loan_is_paid = 0").
		Where("loans.route_id = ?", routeId).
		Where("customers."+searchField+" LIKE ?", "%"+query+"%").
		Find(&loans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load loans",
		})
		return
	}

	if len(loans) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No found",
		})
		return
	}

	var loansResponse []responses.LoanResponse

	for _, loan := range loans {
		loansResponse = append(loansResponse, responses.NewLoanResponse1(loan))
	}

	c.JSON(http.StatusOK, loansResponse)
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

func PayOffLoan(c *gin.Context) {
	id := c.Param("id")

	db := initializers.DB

	var loan models.Loan
	if err := db.Where("id = ?", id).First(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Loan not found"})
		return
	}

	user := c.MustGet("user").(models.User)

	interestAmount := float32(loan.Interest) / 100 * loan.Amount
	fullAmount := loan.Amount + interestAmount

	var fees []models.Fee
	if err := db.Where("loan_id = ?", id).Find(&fees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load fees",
		})
		return
	}

	var message = SetPaymentToPaid(fees, user, fullAmount, interestAmount, db)

	if message != "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": message})
		return
	}

	loan.LoanIsPaid = true

	if err := db.Save(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Loan not saved"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan payment processed successfully"})
}

func SetLoanToPaid(c *gin.Context) {

	id := c.Param("id")

	db := initializers.DB

	var loan models.Loan
	if err := db.Where("id = ?", id).First(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Loan not found"})
		return
	}

	user := c.MustGet("user").(models.User)
	interestAmount := float32(loan.Interest) / 100 * loan.Amount
	fullAmount := loan.Amount + interestAmount

	var fees []models.Fee
	if err := db.Where("loan_id = ?", id).Find(&fees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load fees",
		})
		return
	}

	var message = ""
	message = SetPaymentToPaid(fees, user, fullAmount, interestAmount, db)

	if message != "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": message})
		return
	}

	loan.LoanIsPaid = true

	if err := db.Save(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Loan not saved"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan payment processed successfully"})
}

func CreateRenewLoan(c *gin.Context) {

	var body struct {
		LoanId       string  `form:"loanId" json:"loanId"`
		Amount       float32 `form:"amount" json:"amount"`
		Interest     float32 `form:"interest" json:"interest"`
		FeesQuantity int32   `form:"feesQuantity" json:"feesQuantity"`
	}

	db := initializers.DB

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	//Get the loan
	var loan models.Loan
	if err := db.Where("id = ?", body.LoanId).Find(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load loan",
		})
		return
	}

	var route models.Route
	if err := db.Where("id = ?", loan.RouteId).Find(&route).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load route",
		})
		return
	}

	var user models.User
	if err := initializers.DB.Where("id = ?", loan.UserId).Find(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load user",
		})
		return
	}

	var customer models.Customer
	if err := db.Where("id = ?", loan.CustomerId).Find(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load customer",
		})
		return
	}

	if utils.HaveanactiveloanRenew(customer.ID.String()) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "El prestamo ta pago.",
		})
		return
	}

	var fees []models.Fee
	if err := db.Where("loan_id = ?", body.LoanId).Find(&fees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load fees",
		})
		return
	}

	user1 := c.MustGet("user").(models.User)
	interestAmount := float32(loan.Interest) / 100 * loan.Amount
	fullAmount := loan.Amount + interestAmount

	var message = SetPaymentToPaid(fees, user1, fullAmount, interestAmount, db)

	if message != "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": message})
		return
	}

	oldTotalAmount := 0
	oldInterestAmount := float32(loan.Interest) / 100 * loan.Amount
	oldTotalAmount = int(oldInterestAmount*float32(loan.FeesQuantity)) - oldTotalAmount
	totalAmount := body.Amount + float32(oldTotalAmount)

	loanUUID := uuid.New()
	now := time.Now()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	newLoan := models.Loan{
		Model:         gorm.Model{},
		ID:            loanUUID,
		CustomerId:    customer.ID,
		RouteId:       route.ID,
		UserId:        user.ID,
		Amount:        totalAmount,
		Interest:      body.Interest,
		FeesQuantity:  body.FeesQuantity,
		Date:          currentDate,
		LoanIsPaid:    false,
		IsRenewed:     true,
		IsCurrentLoan: true,
	}

	loan.LoanIsPaid = true
	loan.IsCurrentLoan = false

	db.Save(&loan)

	resultLoan := db.Create(&newLoan)
	if resultLoan.Error != nil {
		fmt.Println(resultLoan.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create loan",
		})
		return
	}

	for i := int32(1); i <= body.FeesQuantity; i++ {
		fee := models.Fee{ID: uuid.New(), LoanId: loanUUID, Number: i, ExpectedDate: currentDate.Add(time.Duration(i) * 7 * 24 * time.Hour)}
		resultFee := db.Create(&fee)
		if resultFee.Error != nil {
			fmt.Println(resultFee.Error)
			initializers.DB.Unscoped().Delete(&loan)
			initializers.DB.Unscoped().Where("loan_id = ?", loanUUID).Delete(models.Fee{})
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create loan",
			})
			return
		}
	}

	var newFees []models.Fee
	if err := db.Where("loan_id = ?", newLoan.ID).Order("number").Find(&newFees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load fee",
		})
		return
	}

	loanRes := responses.NewLoanResponse(newLoan)

	c.JSON(http.StatusCreated, loanRes)
}

func SetPaymentToPaid(fees []models.Fee, user models.User, fullAmount float32, interestAmount float32, db *gorm.DB) string {

	sort.Slice(fees, func(i, j int) bool {
		return fees[i].Number < fees[j].Number
	})

	for _, fee := range fees {
		var payments []models.Payment
		if err := db.Where("fee_id = ?", fee.ID).Find(&payments).Error; err != nil {
			return "Payment not found"
		}

		if len(payments) == 0 {
			// Create a payment with the full amount
			payment := models.Payment{
				ID:         uuid.New(),
				PaidAmount: fullAmount,
				FeeId:      fee.ID,
				PaidDate:   time.Now(),
				UserId:     user.ID,
			}

			if err := db.Create(&payment).Error; err != nil {
				return "Payment not saved"
			}

			return ""
		} else {
			// Check if the fee has been paid
			var paidAmount float32
			for _, payment := range payments {
				paidAmount += payment.PaidAmount
			}

			if paidAmount < interestAmount {
				// Create a payment with the rest of loan amount

				payment := models.Payment{
					ID:         uuid.New(),
					PaidAmount: fullAmount - paidAmount,
					FeeId:      fee.ID,
					PaidDate:   time.Now(),
					UserId:     user.ID,
				}

				if err := db.Create(&payment).Error; err != nil {
					return "Payment not saved"
				}

				return ""
			}
		}
	}

	return ""
}
