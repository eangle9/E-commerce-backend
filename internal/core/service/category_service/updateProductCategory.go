package categoryservice

import (
	"Eccomerce-website/internal/core/common/utils"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productCategoryService) UpdateProductCategory(id int, category utils.UpdateCategory) response.Response {
	productCategory, err := p.categoryRepo.EditProductCategoryById(id, category)
	if err != nil {
		response := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: err.Error(),
		}
		return response
	}

	errorMessage := fmt.Sprintf("you have updated product category with id %d", id)
	response := response.Response{
		Data:         productCategory,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: errorMessage,
	}

	return response
}
