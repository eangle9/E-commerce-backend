package productsservice

import (
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
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       productList,
		StatusCode: http.StatusOK,
		Message:    "you have get all list of products",
	}

	return response
}
