package service

import (
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
)

type ProductCategoryService interface {
	CreateProductCategory(request request.ProductCategoryRequest) response.Response
}
