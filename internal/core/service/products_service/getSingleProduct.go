package productsservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productsService) GetSingleProduct(id int) response.Response {
	singleProduct, err := p.productsRepo.GetSingleProductById(id)
	if err != nil {
		response := response.Response{
			Status:       http.StatusNotFound,
			ErrorType:    errorcode.NotFoundError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         singleProduct,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: fmt.Sprintf("you have get the product with product_id '%d'", id),
	}

	return response
}
