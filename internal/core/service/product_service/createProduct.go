package productservice

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

type productService struct {
	productRepo   repository.ProductRepository
	serviceLogger *zap.Logger
}

func NewProductService(repo repository.ProductRepository, serviceLoggger *zap.Logger) service.ProductService {
	return &productService{
		productRepo:   repo,
		serviceLogger: serviceLoggger,
	}
}

func (p productService) CreateProduct(ctx context.Context, request request.ProductRequest, requestID string) (response.Response, error) {
	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "create-product validation error").WithProperty(entity.StatusCode, 400)
		p.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "CreateProduct"),
			zap.String("requestID", requestID),
			zap.Any("requestData", request),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	product := dto.Product{
		CategoryID:  request.CategoryID,
		Brand:       request.Brand,
		ProductName: request.ProductName,
		Description: request.Description,
	}

	id, err := p.productRepo.InsertProduct(ctx, product, requestID)
	if err != nil {
		return response.Response{}, err
	}

	product.ID = *id

	response := response.Response{
		Data:       product,
		StatusCode: http.StatusCreated,
		Message:    "Congratulation, product created successfully!",
	}

	return response, nil
}
