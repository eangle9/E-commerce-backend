package platform

import (
	// "Eccomerce-website/internal/core/dto"
	// "Eccomerce-website/internals/core/dto"

	"context"
	"mime/multipart"
)

// type API interface {
// 	InitiatePayment(request *dto.PaymentRequest) (*dto.PaymentResponse, error)
// 	VerifyPayment(txRef string) (*dto.VerifyResponse, error)
// }

type Asset interface {
	SaveAsset(ctx context.Context, asset multipart.File, dst string) error
}
