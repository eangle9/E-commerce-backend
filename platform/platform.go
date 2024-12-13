package platform

import "Eccomerce-website/internal/core/dto"

type API interface {
	InitiatePayment(request *dto.PaymentRequest) (*dto.PaymentResponse, error)
	VerifyPayment(txRef string) (*dto.VerifyResponse, error)
}
