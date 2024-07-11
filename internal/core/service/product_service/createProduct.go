package productservice

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
)

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) service.ProductService {
	return &productService{
		productRepo: repo,
	}
}

func (p productService) CreateProduct(request request.ProductRequest) response.Response {
	product := dto.Product{
		CategoryID:  request.CategoryID,
		Brand:       request.Brand,
		ProductName: request.ProductName,
		Description: request.Description,
	}

	id, err := p.productRepo.InsertProduct(product)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusConflict,
			Message:    err.Error(),
		}
		return response
	}

	product.ID = *id

	response := response.Response{
		Data:       product,
		StatusCode: http.StatusCreated,
		Message:    "Congratulation, product created successfully!",
	}

	return response
}
