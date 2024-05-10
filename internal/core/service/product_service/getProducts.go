package productservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (p productService) GetProducts() response.Response {
	products, err := p.productRepo.ListProducts()
	if err != nil {
		response := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         products,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: "you have get all list of products",
	}

	return response
}
