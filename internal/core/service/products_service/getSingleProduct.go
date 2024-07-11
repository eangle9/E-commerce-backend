package productsservice

import (
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productsService) GetSingleProduct(id int) response.Response {
	singleProduct, err := p.productsRepo.GetSingleProductById(id)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       singleProduct,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get the product with product_id '%d'", id),
	}

	return response
}
