package responses

import "prestamosbackend/models"

type SignatureTypeResponse struct {
	ID   string `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

func NewSignatureTypeResponse(signatureType models.SignatureType) *SignatureTypeResponse {
	return &SignatureTypeResponse{
		ID:   signatureType.ID.String(),
		Name: signatureType.Name,
	}
}
