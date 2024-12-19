package categoryservice

import (
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (p productCategoryService) GetProductCategory(ctx context.Context, id int, requestID string) (response.Response, error) {
	category, err := p.categoryRepo.GetProductCategoryById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       category,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get product category with category_id %d", id),
	}

	return response, nil
}
