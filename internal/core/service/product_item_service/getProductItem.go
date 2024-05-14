package productitemservice

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productItemService) GetProductItem(id int) response.Response {
	productItem, err := p.itemRepo.GetProductItemById(id)
	if err != nil {
		response := response.Response{
			Status:       http.StatusNotFound,
			ErrorType:    errorcode.NotFoundError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         productItem,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: fmt.Sprintf("you have get single product item with id '%d'", id),
	}

	return response
}
