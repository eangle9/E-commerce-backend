package productitemservice

import (
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (p productItemService) DeleteProductItem(id int) response.Response {
	resp, status, _, err := p.itemRepo.DeleteProductItemById(id)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       resp,
		StatusCode: status,
		Message:    "you have get all list of colors",
	}

	return response
}
