package productservice

import (
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (p productService) DeleteProduct(ctx context.Context, id int, requestID string) (response.Response, error) {
	err := p.productRepo.DeleteProductById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("product with product_id %d deleted successfully", id),
	}

	return response, nil
}
