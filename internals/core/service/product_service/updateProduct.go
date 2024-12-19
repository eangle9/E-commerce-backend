package productservice

import (
	"Eccomerce-website/internals/core/common/utils"
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (p productService) UpdateProduct(ctx context.Context, id int, product utils.UpdateProduct, requestID string) (response.Response, error) {
	updatedProduct, err := p.productRepo.EditProductById(ctx, id, product, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       updatedProduct,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have successfully updated product with product_id '%d'", id),
	}

	return response, nil
}
