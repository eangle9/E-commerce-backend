package productitemservice

import (
	"Eccomerce-website/internal/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (p productItemService) GetProductItem(ctx context.Context, id int, requestID string) (response.Response, error) {
	productItem, err := p.itemRepo.GetProductItemById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       productItem,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get single product item with id '%d'", id),
	}

	return response, nil
}
