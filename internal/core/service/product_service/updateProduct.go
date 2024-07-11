package productservice

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productService) UpdateProduct(id int, product utils.UpdateProduct) response.Response {
	updatedProduct, err := p.productRepo.EditProductById(id, product)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       updatedProduct,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get product with product_id '%d'", id),
	}

	return response
}
