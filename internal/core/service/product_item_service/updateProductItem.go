package productitemservice

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productItemService) UpdateProductItem(id int, item utils.UpdateProductItem) response.Response {
	productItem, err := p.itemRepo.EditProductItemById(id, item)
	if err != nil {
		response := response.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:       productItem,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have successfully updated product item with product_item_id '%d'", id),
	}

	return response
}
