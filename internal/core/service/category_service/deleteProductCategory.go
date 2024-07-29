package categoryservice

import (
	"Eccomerce-website/internal/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (p productCategoryService) DeleteProductCategory(ctx context.Context, id int, requestID string) (response.Response, error) {
	err := p.categoryRepo.DeleteProductCategoryById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("category with category_id %d deleted successfully", id),
	}

	return response, nil
}
