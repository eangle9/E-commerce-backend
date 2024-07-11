package categoryservice

import (
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (p productCategoryService) GetProductCategories() response.Response {
	productCategories, err := p.categoryRepo.ListProductCategory()
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusInternalServerError,
			// ErrorMessage: "failed to get list of product categories",
			Message: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       productCategories,
		StatusCode: http.StatusOK,
		Message:    "you have get list of product categories successfully!",
	}

	return response
}
