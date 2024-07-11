package productservice

import (
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productService) GetProducts() response.Response {
	products, err := p.productRepo.ListProducts()
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       products,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get product with product_id "),
	}

	return response
}
