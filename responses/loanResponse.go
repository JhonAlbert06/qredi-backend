package responses

import (
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"sort"
)

type LoanResponse struct {
	ID            string                 `json:"id" form:"id"`
	Amount        float32                `json:"amount" form:"amount"`
	Interest      float32                `json:"interest" form:"interest"`
	FeesQuantity  int32                  `json:"feesQuantity" form:"feesQuantity"`
	LoanIsPaid    bool                   `json:"loanIsPaid" form:"loanIsPaid"`
	IsRenewed     bool                   `json:"isRenewed" form:"isRenewed"`
	IsCurrentLoan bool                   `json:"isCurrentLoan" form:"isCurrentLoan"`
	Date          models.Date            `json:"date" form:"date"`
	Customer      *CustomerResponse      `json:"customer" form:"customer"`
	Route         *RouteResponse         `json:"route" form:"route"`
	User          *UserResponse          `json:"user" form:"user"`
	Fees          []*FeeResponse         `json:"fees" form:"fees"`
}

type LoanConfirmationResponse struct {
	Customer           *CustomerResponse `json:"customer" form:"customer"`
	Route              *RouteResponse    `json:"route" form:"route"`
	Amount             float32           `json:"amount" form:"amount"`
	InterestAmount     float32           `json:"interestAmount" form:"interestAmount"`
	InterestPercentage float32           `json:"interestPercentage" form:"interestPercentage"`
	FeesQuantity       int32             `json:"feesQuantity" form:"feesQuantity"`
	DateCreation       models.Date       `json:"dateCreation" form:"dateCreation"`
	DateFirst          models.Date       `json:"dateFirst" form:"dateFirst"`
	DateLast           models.Date       `json:"dateLast" form:"dateLast"`
	AmountFinal        float32           `json:"amountFinal" form:"amountFinal"`
	ProfitsAmount      float32           `json:"profitsAmount" form:"profitsAmount"`
	ProfitsPercentage  float32           `json:"profitsPercentage" form:"profitsPercentage"`
}

type ReNewLoanConfirmationResponse struct {
	Customer           *CustomerResponse `json:"customer" form:"customer"`
	Route              *RouteResponse    `json:"route" form:"route"`
	OldInterestAmount  float32           `json:"oldInterestAmount" form:"oldInterestAmount"`
	OldTotalAmount     int               `json:"oldTotalAmount" form:"oldTotalAmount"`
	Amount             float32           `json:"amount" form:"amount"`
	TotalAmount        float32           `json:"totalAmount" form:"totalAmount"`
	InterestAmount     float32           `json:"interestAmount" form:"interestAmount"`
	InterestPercentage float32           `json:"interestPercentage" form:"interestPercentage"`
	FeesQuantity       int32             `json:"feesQuantity" form:"feesQuantity"`
	DateCreation       models.Date       `json:"dateCreation" form:"dateCreation"`
	DateFirst          models.Date       `json:"dateFirst" form:"dateFirst"`
	DateLast           models.Date       `json:"dateLast" form:"dateLast"`
	AmountFinal        float32           `json:"amountFinal" form:"amountFinal"`
	ProfitsAmount      float32           `json:"profitsAmount" form:"profitsAmount"`
	ProfitsPercentage  float32           `json:"profitsPercentage" form:"profitsPercentage"`
}

func NewLoanResponse(loan models.Loan) *LoanResponse {
	db := initializers.DB

	//var date models.Date
	date := models.ToDate(loan.Date)

	var route models.Route
	if err := db.First(&route, loan.RouteId).Error; err != nil {
		route = models.Route{}
	}

	var customer models.Customer
	if err := db.First(&customer, loan.CustomerId).Error; err != nil {
		customer = models.Customer{}
	}

	var user models.User
	if err := db.First(&user, loan.UserId).Error; err != nil {
		user = models.User{}
	}

	var fees []models.Fee
	if err := db.Where("loan_id = ?", loan.ID).Find(&fees).Error; err != nil {
		fees = []models.Fee{}
	}

	sort.Slice(fees, func(i, j int) bool {
		return fees[i].Number < fees[j].Number
	})

	var feesResponse []*FeeResponse
	for _, fee := range fees {
		feesResponse = append(feesResponse, NewFeeResponse(fee))
	}

	return &LoanResponse{
		ID:            loan.ID.String(),
		Amount:        loan.Amount,
		Interest:      loan.Interest,
		FeesQuantity:  loan.FeesQuantity,
		LoanIsPaid:    loan.LoanIsPaid,
		IsRenewed:     loan.IsRenewed,
		IsCurrentLoan: loan.IsCurrentLoan,
		Date:          date,
		Customer:      NewCustomerResponse(customer),
		Route:         NewRouteResponse(route),
		User:          NewUserResponse(user),
		Fees:          feesResponse,
	}
}

func NewLoanResponse1(loan models.Loan) LoanResponse {
	db := initializers.DB

	//var date models.Date
	date := models.ToDate(loan.Date)

	var route models.Route
	if err := db.First(&route, loan.RouteId).Error; err != nil {
		route = models.Route{}
	}

	var customer models.Customer
	if err := db.First(&customer, loan.CustomerId).Error; err != nil {
		customer = models.Customer{}
	}

	var user models.User
	if err := db.First(&user, loan.UserId).Error; err != nil {
		user = models.User{}
	}

	var fees []models.Fee
	if err := db.Where("loan_id = ?", loan.ID).Find(&fees).Error; err != nil {
		fees = []models.Fee{}
	}

	sort.Slice(fees, func(i, j int) bool {
		return fees[i].Number < fees[j].Number
	})

	var feesResponse []*FeeResponse
	for _, fee := range fees {
		feesResponse = append(feesResponse, NewFeeResponse(fee))
	}

	return LoanResponse{
		ID:            loan.ID.String(),
		Amount:        loan.Amount,
		Interest:      loan.Interest,
		FeesQuantity:  loan.FeesQuantity,
		LoanIsPaid:    loan.LoanIsPaid,
		IsRenewed:     loan.IsRenewed,
		IsCurrentLoan: loan.IsCurrentLoan,
		Date:          date,
		Customer:      NewCustomerResponse(customer),
		Route:         NewRouteResponse(route),
		User:          NewUserResponse(user),
		Fees:          feesResponse,
	}
}

func NewLoanConfirmationResponse(loan models.LoanConfirmation) *LoanConfirmationResponse {
	db := initializers.DB

	var route models.Route
	if err := db.First(&route, loan.RouteId).Error; err != nil {
		return nil
	}

	var customer models.Customer
	if err := db.First(&customer, loan.CustomerId).Error; err != nil {
		return nil
	}

	return &LoanConfirmationResponse{
		Customer:           NewCustomerResponse(customer),
		Route:              NewRouteResponse(route),
		Amount:             loan.Amount,
		InterestAmount:     loan.InterestAmount,
		InterestPercentage: loan.InterestPercentage,
		FeesQuantity:       loan.FeesQuantity,
		DateCreation:       models.ToDate(loan.DateCreation),
		DateFirst:          models.ToDate(loan.DateFirst),
		DateLast:           models.ToDate(loan.DateLast),
		AmountFinal:        loan.AmountFinal,
		ProfitsAmount:      loan.ProfitsAmount,
		ProfitsPercentage:  loan.ProfitsPercentage,
	}
}

func NewReNewLoanConfirmationResponse(loan models.ReNewLoanConfirmation) *ReNewLoanConfirmationResponse {
	db := initializers.DB

	var route models.Route
	if err := db.First(&route, loan.RouteId).Error; err != nil {
		return nil
	}

	var customer models.Customer
	if err := db.First(&customer, loan.CustomerId).Error; err != nil {
		return nil
	}

	return &ReNewLoanConfirmationResponse{
		Customer:           NewCustomerResponse(customer),
		Route:              NewRouteResponse(route),
		OldInterestAmount:  loan.OldInterestAmount,
		OldTotalAmount:     loan.OldTotalAmount,
		Amount:             loan.Amount,
		TotalAmount:        loan.TotalAmount,
		InterestAmount:     loan.InterestAmount,
		InterestPercentage: loan.InterestPercentage,
		FeesQuantity:       loan.FeesQuantity,
		DateCreation:       models.ToDate(loan.DateCreation),
		DateFirst:          models.ToDate(loan.DateFirst),
		DateLast:           models.ToDate(loan.DateLast),
		AmountFinal:        loan.AmountFinal,
		ProfitsAmount:      loan.ProfitsAmount,
		ProfitsPercentage:  loan.ProfitsPercentage,
	}
}
