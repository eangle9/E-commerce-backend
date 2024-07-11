package categoryservice

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productCategoryService) UpdateProductCategory(id int, category utils.UpdateCategory) response.Response {
	productCategory, err := p.categoryRepo.EditProductCategoryById(id, category)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}
		return response
	}

	errorMessage := fmt.Sprintf("you have updated product category with id %d", id)
	response := response.Response{
		Data:       productCategory,
		StatusCode: http.StatusOK,
		Message:    errorMessage,
	}

	return response
}
