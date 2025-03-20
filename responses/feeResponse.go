package responses

import (
	"prestamosbackend/initializers"
	"prestamosbackend/models"
)

type FeeResponse struct {
	ID            string                `json:"id" form:"id"`
	Number        int32                 `json:"number" form:"number"`
	ExpectedDate  models.Date           `json:"expectedDate" form:"expectedDate"`
	Payments      []PaymentResponse     `json:"payments" form:"payments"`
}

func NewFeeResponse(fee models.Fee) *FeeResponse {
	db := initializers.DB

	var signatureType models.SignatureType
	if err := db.First(&signatureType, fee.ID).Error; err != nil {
		signatureType = models.SignatureType{}
	}

	var payments []models.Payment
	if err := db.Where("fee_id = ?", fee.ID).Find(&payments).Error; err != nil {
		payments = []models.Payment{}
	}

	var paymentsResponse []PaymentResponse
	for _, payment := range payments {
		paymentsResponse = append(paymentsResponse, *NewPaymentResponse(payment))
	}

	return &FeeResponse{
		ID:            fee.ID.String(),
		Number:        fee.Number,
		ExpectedDate:  models.ToDate(fee.ExpectedDate),
		Payments:      paymentsResponse,
	}
}
