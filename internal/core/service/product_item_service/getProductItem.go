package productitemservice

import (
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productItemService) GetProductItem(id int) response.Response {
	productItem, err := p.itemRepo.GetProductItemById(id)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       productItem,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get single product item with id '%d'", id),
	}

	return response
}
