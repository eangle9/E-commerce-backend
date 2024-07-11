package productservice

import (
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productService) GetProduct(id int) response.Response {
	product, err := p.productRepo.GetProductById(id)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       product,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get product with product_id '%d'", id),
	}

	return response
}
