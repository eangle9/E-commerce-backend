package categoryservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (p productCategoryService) GetProductCategories() response.Response {
	productCategories, err := p.categoryRepo.ListProductCategory()
	if err != nil {
		response := response.Response{
			Status:    http.StatusInternalServerError,
			ErrorType: errorcode.InternalError,
			// ErrorMessage: "failed to get list of product categories",
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         productCategories,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: "you have get list of product categories successfully!",
	}

	return response
}
