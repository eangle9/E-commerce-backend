package colorservice

import (
	"Eccomerce-website/internal/core/dto"
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
			StatusCode: http.StatusConflict,
			Message:    err.Error(),
		}
		return response
	}

	color.ID = *id

	response := response.Response{
		Data:       color,
		StatusCode: http.StatusCreated,
		Message:    "Congratulation, color for product created successfully!",
	}

	return response
}
