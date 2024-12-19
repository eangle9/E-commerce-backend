package productitemservice

import (
	"Eccomerce-website/internals/core/common/utils"
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (p productItemService) UpdateProductItem(ctx context.Context, id int, item utils.UpdateProductItem, requestID string) (response.Response, error) {
	productItem, err := p.itemRepo.EditProductItemById(ctx, id, item, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       productItem,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have successfully updated product item with product_item_id '%d'", id),
	}

	return response, nil
}
