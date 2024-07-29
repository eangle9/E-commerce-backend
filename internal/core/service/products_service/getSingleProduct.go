package productsservice

import (
	"Eccomerce-website/internal/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (p productsService) GetSingleProduct(ctx context.Context, id int, requestID string) (response.Response, error) {
	singleProduct, err := p.productsRepo.GetSingleProductById(ctx, id, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       singleProduct,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have get the product with product_id '%d'", id),
	}

	return response, nil
}
