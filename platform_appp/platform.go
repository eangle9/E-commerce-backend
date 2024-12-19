package platform

import "Eccomerce-website/internals/core/dto"

type API interface {
	InitiatePayment(request *dto.PaymentRequest) (*dto.PaymentResponse, error)
	VerifyPayment(txRef string) (*dto.VerifyResponse, error)
}
