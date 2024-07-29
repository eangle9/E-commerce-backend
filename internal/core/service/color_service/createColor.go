package colorservice

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type colorService struct {
	colorRepo     repository.ColorRepository
	serviceLogger *zap.Logger
}

func NewColorService(colorRepo repository.ColorRepository, serviceLogger *zap.Logger) service.ColorService {
	return &colorService{
		colorRepo:     colorRepo,
		serviceLogger: serviceLogger,
	}
}

func (c colorService) CreateColor(ctx context.Context, request request.ColorRequest, requestID string) (response.Response, error) {
	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "create-color validation error").WithProperty(entity.StatusCode, 400)
		c.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "CreateColor"),
			zap.String("requestID", requestID),
			zap.Any("requestData", request),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	color := dto.Color{
		Name: request.ColorName,
	}

	id, err := c.colorRepo.InsertColor(ctx, color, requestID)
	if err != nil {
		return response.Response{}, err
	}

	color.ID = *id

	response := response.Response{
		Data:       color,
		StatusCode: http.StatusCreated,
		Message:    "Congratulation, color for product created successfully!",
	}

	return response, nil
}
