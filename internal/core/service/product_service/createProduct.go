package productservice

import (
	"Eccomerce-website/internal/core/dto"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
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
		ProductName: request.ProductName,
		Description: request.Description,
	}

	id, err := p.productRepo.InsertProduct(product)
	if err != nil {
		response := response.Response{
			Status:       http.StatusConflict,
			ErrorType:    "DUPLICATE_ENTRY",
			ErrorMessage: err.Error(),
		}
		return response
	}

	product.ID = *id

	response := response.Response{
		Data:         product,
		Status:       http.StatusCreated,
		ErrorType:    errorcode.Success,
		ErrorMessage: "Congratulation, product created successfully!",
	}

	return response
}
