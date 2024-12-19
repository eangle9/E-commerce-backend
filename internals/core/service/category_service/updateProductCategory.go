package categoryservice

import (
	"Eccomerce-website/internals/core/common/utils"
	"Eccomerce-website/internals/core/model/response"
	"context"
	"fmt"
	"net/http"
)

func (p productCategoryService) UpdateProductCategory(ctx context.Context, id int, category utils.UpdateCategory, requestID string) (response.Response, error) {
	productCategory, err := p.categoryRepo.EditProductCategoryById(ctx, id, category, requestID)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       productCategory,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have updated product category with id %d", id),
	}

	return response, nil
}
