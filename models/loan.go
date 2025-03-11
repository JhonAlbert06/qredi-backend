package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Loan struct {
	ID            uuid.UUID `gorm:"type:char(36);primary_key;"`
	CustomerId    uuid.UUID `gorm:"type:char(36)"`
	RouteId       uuid.UUID `gorm:"type:char(36)"`
	UserId        uuid.UUID `gorm:"type:char(36)"`
	Amount        float32
	Interest      float32
	FeesQuantity  int32
	LoanIsPaid    bool
	IsRenewed     bool
	IsCurrentLoan bool
	Date          time.Time
	gorm.Model
}

type LoanConfirmation struct {
	CustomerId         uuid.UUID `json:"customerId"`
	RouteId            uuid.UUID `json:"routeId"`
	Amount             float32   `json:"amount"`
	InterestAmount     float32   `json:"interestAmount"`
	InterestPercentage float32   `json:"interestPercentage"`
	FeesQuantity       int32     `json:"feesQuantity"`
	DateCreation       time.Time `json:"dateCreation"`
	DateFirst          time.Time `json:"dateFirst"`
	DateLast           time.Time `json:"dateLast"`
	AmountFinal        float32   `json:"amountFinal"`
	ProfitsAmount      float32   `json:"profitsAmount"`
	ProfitsPercentage  float32   `json:"profitsPercentage"`
}

type ReNewLoanConfirmation struct {
	CustomerId         uuid.UUID `json:"customerId"`
	RouteId            uuid.UUID `json:"routeId"`
	OldInterestAmount  float32   `json:"oldInterestAmount" form:"oldInterestAmount"`
	OldTotalAmount     int       `json:"oldTotalAmount" form:"oldTotalAmount"`
	Amount             float32   `json:"amount" form:"amount"`
	TotalAmount        float32   `json:"totalAmount" form:"totalAmount"`
	InterestAmount     float32   `json:"interestAmount" form:"interestAmount"`
	InterestPercentage float32   `json:"interestPercentage" form:"interestPercentage"`
	FeesQuantity       int32     `json:"feesQuantity" form:"feesQuantity"`
	DateCreation       time.Time `json:"dateCreation"`
	DateFirst          time.Time `json:"dateFirst"`
	DateLast           time.Time `json:"dateLast"`
	AmountFinal        float32   `json:"amountFinal" form:"amountFinal"`
	ProfitsAmount      float32   `json:"profitsAmount" form:"profitsAmount"`
	ProfitsPercentage  float32   `json:"profitsPercentage" form:"profitsPercentage"`
}
