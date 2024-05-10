package colorservice

import (
	"Eccomerce-website/internal/core/dto"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
)

type colorService struct {
	colorRepo repository.ColorRepository
}

func NewColorService(colorRepo repository.ColorRepository) service.ColorService {
	return &colorService{
		colorRepo: colorRepo,
	}
}

func (c colorService) CreateColor(request request.ColorRequest) response.Response {
	color := dto.Color{
		Name: request.ColorName,
	}

	id, err := c.colorRepo.InsertColor(color)
	if err != nil {
		response := response.Response{
			Status:       http.StatusConflict,
			ErrorType:    "DUPLICATE_ENTRY",
			ErrorMessage: err.Error(),
		}
		return response
	}

	color.ID = *id

	response := response.Response{
		Data:         color,
		Status:       http.StatusCreated,
		ErrorType:    errorcode.Success,
		ErrorMessage: "you have successfully created color for the product!",
	}

	return response
}
