package sizeservice

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
	"strings"
)

type sizeService struct {
	sizeRepo repository.SizeRepository
}

func NewSizeService(sizeRepo repository.SizeRepository) service.SizeService {
	return &sizeService{
		sizeRepo: sizeRepo,
	}
}

func (s sizeService) CreateSize(request request.SizeRequest) response.Response {
	request.SizeName = strings.ToUpper(request.SizeName)
	size := dto.Size{
		ProductItemID: request.ProductItemID,
		SizeName:      request.SizeName,
		Price:         request.Price,
		Discount:      request.Discount,
		QtyInStock:    request.QtyInStock,
	}

	id, err := s.sizeRepo.InsertSize(size)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusConflict,
			Message:    err.Error(),
		}
		return response
	}

	size.ID = *id

	response := response.Response{
		Data:       size,
		StatusCode: http.StatusCreated,
		Message:    "Congratulation, size for product created successfully!",
	}

	return response
}
