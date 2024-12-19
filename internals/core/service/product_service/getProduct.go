package productservice

import (
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (p productService) GetProduct(ctx context.Context, id int, requestID string) (response.Response, error) {
	product, err := p.productRepo.GetProductById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       product,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get product with product_id '%d'", id),
	}

	return response, nil
}
