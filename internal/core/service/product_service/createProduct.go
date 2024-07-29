package productservice

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"context"
	"net/http"

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
