package productservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productService) GetProduct(id int) response.Response {
	product, err := p.productRepo.GetProductById(id)
	if err != nil {
		response := response.Response{
			Status:       http.StatusNotFound,
			ErrorType:    errorcode.NotFoundError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         product,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: fmt.Sprintf("you have get product with product_id '%d'", id),
	}

	return response
}
