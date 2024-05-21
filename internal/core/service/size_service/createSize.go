package sizeservice

import (
	"Eccomerce-website/internal/core/dto"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
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
	size := dto.Size{
		SizeName: request.SizeName,
	}

	id, err := s.sizeRepo.InsertSize(size)
	if err != nil {
		response := response.Response{
			Status:       http.StatusConflict,
			ErrorType:    "DUPLICATE_ENTRY",
			ErrorMessage: err.Error(),
		}
		return response
	}

	size.ID = *id

	response := response.Response{
		Data:         size,
		Status:       http.StatusCreated,
		ErrorType:    errorcode.Success,
		ErrorMessage: "Congratulation, size for product created successfully!",
	}

	return response
}
