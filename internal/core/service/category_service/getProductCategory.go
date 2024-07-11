package categoryservice

import (
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productCategoryService) GetProductCategory(id int) response.Response {
	category, err := p.categoryRepo.GetProductCategoryById(id)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
		return response
	}

	errorMessage := fmt.Sprintf("you have get product category with category_id %d", id)
	response := response.Response{
		Data:       category,
		StatusCode: http.StatusOK,
		Message:    errorMessage,
	}

	return response
}
