package productitemservice

import (
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productItemService) GetProductItems() response.Response {
	productItems, err := p.itemRepo.ListProductItems()
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       productItems,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get single product item "),
	}

	return response
}
