package productservice

import (
	"Eccomerce-website/internal/core/common/utils"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productService) UpdateProduct(id int, product utils.UpdateProduct) response.Response {
	updatedProduct, err := p.productRepo.EditProductById(id, product)
	if err != nil {
		response := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         updatedProduct,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: fmt.Sprintf("you have successfully updated product with product_id '%d'", id),
	}

	return response
}
