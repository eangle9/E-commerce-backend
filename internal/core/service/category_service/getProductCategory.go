package categoryservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productCategoryService) GetProductCategory(id int) response.Response {
	category, err := p.categoryRepo.GetProductCategoryById(id)
	if err != nil {
		response := response.Response{
			Status:       http.StatusNotFound,
			ErrorType:    errorcode.NotFoundError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	errorMessage := fmt.Sprintf("you have get product category with category_id %d", id)
	response := response.Response{
		Data:         category,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: errorMessage,
	}

	return response
}
