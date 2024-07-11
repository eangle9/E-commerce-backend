package productimageservice

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
)

type productImageService struct {
	imageRepo repository.ProductImageRepository
}

func NewProductImageService(imageRepo repository.ProductImageRepository) service.ProductImageService {
	return &productImageService{
		imageRepo: imageRepo,
	}
}

func (p productImageService) CreateProductImage(request request.ProductImageRequest) response.Response {
	id, imageUrl, err := p.imageRepo.InsertProductImage(request)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return response
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

	return response
}
