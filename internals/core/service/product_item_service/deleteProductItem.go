package productitemservice

import (
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (p productItemService) DeleteProductItem(ctx context.Context, id int, requestID string) (response.Response, error) {
	err := p.itemRepo.DeleteProductItemById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("productItem with product_item_id %d deleted successfully", id),
	}

	return response, nil
}
