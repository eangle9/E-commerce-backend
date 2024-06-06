package productsservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
)

type productsService struct {
	productsRepo repository.GetProducts
}

func NewProductsService(productsRepo repository.GetProducts) service.GetProductService {
	return &productsService{
		productsRepo: productsRepo,
	}
}

func (p productsService) GetAllProducts() response.Response {
	productList, err := p.productsRepo.ListAllProducts()
	if err != nil {
		response := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         productList,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: "you have get all list of products",
	}

	return response
}
