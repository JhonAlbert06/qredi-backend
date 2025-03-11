package responses

import (
	"prestamosbackend/initializers"
	"prestamosbackend/models"
)

type PaymentResponse struct {
	ID         string        `json:"id" form:"id"`
	PaidAmount float32       `json:"paidAmount" form:"paidAmount"`
	PaidDate   models.Date   `json:"paidDate" form:"paidDate"`
	User       *UserResponse `json:"user" form:"user"`
}

func NewPaymentResponse(payment models.Payment) *PaymentResponse {

	db := initializers.DB
	var user models.User
	if err := db.First(&user, payment.UserId).Error; err != nil {
		return &PaymentResponse{}
	}

	return &PaymentResponse{
		ID:         payment.ID.String(),
		PaidAmount: payment.PaidAmount,
		PaidDate:   models.ToDate(payment.PaidDate),
		User:       NewUserResponse(user),
	}
}
