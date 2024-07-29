package service

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"context"
)

type UserService interface {
	SignUp(ctx context.Context, request request.SignUpRequest, requestID string) (response.Response, error)
	LoginUser(ctx context.Context, request request.LoginRequest, requestID string) (response.Response, error)
	GetUsers(ctx context.Context, request request.PaginationQuery, requestID string) (response.Response, error)
	GetUser(ctx context.Context, id int, requestID string) (response.Response, error)
	UpdateUser(ctx context.Context, id int, user request.UpdateUser, requestID string) (response.Response, error)
	DeleteUser(ctx context.Context, id int, requestID string) (response.Response, error)
	RefreshToken(ctx context.Context, refreshToken request.RefreshRequest, requestID string) (response.Response, error)
}

type ProductCategoryService interface {
	CreateProductCategory(request request.ProductCategoryRequest) response.Response
	GetProductCategories() response.Response
	GetProductCategory(id int) response.Response
	UpdateProductCategory(id int, category utils.UpdateCategory) response.Response
	DeleteProductCategory(id int) response.Response
}

type ColorService interface {
	CreateColor(request request.ColorRequest) response.Response
	GetColors() response.Response
	GetColor(id int) response.Response
	UpdateColor(id int, color utils.UpdateColor) response.Response
	DeleteColor(id int) response.Response
}

type ProductService interface {
	CreateProduct(ctx context.Context, request request.ProductRequest, requestID string) (response.Response, error)
	GetProducts(ctx context.Context, request request.PaginationQuery, requestID string) (response.Response, error)
	GetProduct(ctx context.Context, id int, requestID string) (response.Response, error)
	UpdateProduct(ctx context.Context, id int, product utils.UpdateProduct, requestID string) (response.Response, error)
	DeleteProduct(ctx context.Context, id int, requestID string) (response.Response, error)
}

type ProductItemService interface {
	CreateProductItem(ctx context.Context, request request.ProductItemRequest, requestID string) (response.Response, error)
	GetProductItems(ctx context.Context, request request.PaginationQuery, requestID string) (response.Response, error)
	GetProductItem(ctx context.Context, id int, requestID string) (response.Response, error)
	UpdateProductItem(ctx context.Context, id int, productItem utils.UpdateProductItem, requestID string) (response.Response, error)
	DeleteProductItem(ctx context.Context, id int, requestID string) (response.Response, error)
}

type ProductImageService interface {
	CreateProductImage(request request.ProductImageRequest) response.Response
}

type CartService interface {
	AddToCart(request request.CartRequest, userId uint) response.Response
}

type SizeService interface {
	CreateSize(ctx context.Context, request request.SizeRequest, requestID string) (response.Response, error)
	GetSizes(ctx context.Context, request request.PaginationQuery, requestID string) (response.Response, error)
	GetSize(ctx context.Context, id int, requestID string) (response.Response, error)
	UpdateSize(ctx context.Context, id int, size request.UpdateSize, requestID string) (response.Response, error)
	DeleteSize(ctx context.Context, id int, requestID string) (response.Response, error)
}

type GetProductService interface {
	GetAllProducts(ctx context.Context, request request.PaginationQuery, requestID string) (response.Response, error)
	GetSingleProduct(ctx context.Context, id int, requestID string) (response.Response, error)
}

type ReviewService interface {
	CreateReview(ctx context.Context, request request.ReviewRequest, requestID string) (response.Response, error)
	GetReviews(ctx context.Context, paginationQuery request.PaginationQuery, requestID string) (response.Response, error)
}
