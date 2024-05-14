package productitemservice

import (
	"Eccomerce-website/internal/core/common/utils"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (p productItemService) UpdateProductItem(id int, item utils.UpdateProductItem) response.Response {
	productItem, err := p.itemRepo.EditProductItemById(id, item)
	if err != nil {
		response := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.InvalidRequest,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Data:         productItem,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: fmt.Sprintf("you have successfully updated product item with product_item_id '%d'", id),
	}

	return response
}
