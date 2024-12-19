package productimageservice

import (
	"Eccomerce-website/internals/core/dto"
	"Eccomerce-website/internals/core/model/request"
	"Eccomerce-website/internals/core/model/response"
	"Eccomerce-website/internals/core/port/repository"
	"Eccomerce-website/internals/core/port/service"
	"context"
	"net/http"

	"go.uber.org/zap"
)

type productImageService struct {
	imageRepo     repository.ProductImageRepository
	serviceLogger *zap.Logger
}

func NewProductImageService(imageRepo repository.ProductImageRepository, serviceLogger *zap.Logger) service.ProductImageService {
	return &productImageService{
		imageRepo:     imageRepo,
		serviceLogger: serviceLogger,
	}
}

func (p productImageService) CreateProductImage(ctx context.Context, request request.ProductImageRequest, requestID string) (response.Response, error) {
	id, imageUrl, err := p.imageRepo.InsertProductImage(ctx, request, requestID)
	if err != nil {
		return response.Response{}, err
	}

	productImage := dto.ProductImage{
		ID:            *id,
		ProductItemID: request.ProductItemId,
		ImageUrl:      imageUrl,
	}

	response := response.Response{
		Data:       productImage,
		StatusCode: http.StatusOK,
		Message:    "upload successful!",
	}

	return response, nil
}
