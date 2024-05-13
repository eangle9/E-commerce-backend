package productitemservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (p productItemService) GetProductItems() response.Response {
	productItems, err := p.itemRepo.ListProductItems()
	if err != nil {
		response := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         productItems,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: "you have get all list of product items!",
	}

	return response
}
