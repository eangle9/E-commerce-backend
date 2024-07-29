package categoryservice

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

type productCategoryService struct {
	categoryRepo  repository.ProductCategoryRepository
	serviceLogger *zap.Logger
}

func NewProductCategoryRepository(repo repository.ProductCategoryRepository, serviceLogger *zap.Logger) service.ProductCategoryService {
	return &productCategoryService{
		categoryRepo:  repo,
		serviceLogger: serviceLogger,
	}
}

func (p productCategoryService) CreateProductCategory(ctx context.Context, request request.ProductCategoryRequest, requestID string) (response.Response, error) {
	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "create-category validation error").WithProperty(entity.StatusCode, 400)
		p.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "CreateProductCategory"),
			zap.String("requestID", requestID),
			zap.Any("requestData", request),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	parentId := request.ParentID
	name := request.Name

	category := dto.ProductCategory{
		ParentID: parentId,
		Name:     name,
	}

	id, err := p.categoryRepo.InsertProductCategory(ctx, category, requestID)
	if err != nil {
		return response.Response{}, err
	}

	category.ID = *id

	response := response.Response{
		Data:       category,
		StatusCode: http.StatusCreated,
		Message:    "Congratulation, product category created successfully!",
	}

	return response, nil
}
