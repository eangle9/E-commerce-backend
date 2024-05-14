package productimageservice

import (
	"Eccomerce-website/internal/core/dto"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
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
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	productImage := dto.ProductImage{
		ID:            *id,
		ProductItemID: request.ProductItemId,
		ImageUrl:      imageUrl,
	}

	response := response.Response{
		Data:         productImage,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: "upload successful!",
	}

	return response
}
