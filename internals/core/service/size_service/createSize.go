package sizeservice

import (
	"Eccomerce-website/internals/core/dto"
	"Eccomerce-website/internals/core/entity"
	"Eccomerce-website/internals/core/model/request"
	"Eccomerce-website/internals/core/model/response"
	"Eccomerce-website/internals/core/port/repository"
	"Eccomerce-website/internals/core/port/service"
	"context"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type sizeService struct {
	sizeRepo      repository.SizeRepository
	serviceLogger *zap.Logger
}

func NewSizeService(sizeRepo repository.SizeRepository, serviceLogger *zap.Logger) service.SizeService {
	return &sizeService{
		sizeRepo:      sizeRepo,
		serviceLogger: serviceLogger,
	}
}

func (s sizeService) CreateSize(ctx context.Context, request request.SizeRequest, requestID string) (response.Response, error) {
	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "createSize validation error").WithProperty(entity.StatusCode, 400)
		s.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "CreateSize"),
			zap.String("requestID", requestID),
			zap.Any("requestData", request),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	request.SizeName = strings.ToUpper(request.SizeName)
	size := dto.Size{
		ProductItemID: request.ProductItemID,
		SizeName:      request.SizeName,
		Price:         request.Price,
		Discount:      request.Discount,
		QtyInStock:    request.QtyInStock,
	}

	id, err := s.sizeRepo.InsertSize(ctx, size, requestID)
	if err != nil {
		return response.Response{}, err
	}

	size.ID = *id

	response := response.Response{
		Data:       size,
		StatusCode: http.StatusCreated,
		Message:    "product size created successfully!",
	}

	return response, nil
}
