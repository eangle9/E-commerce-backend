package service

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
)

type UserService interface {
	SignUp(request request.SignUpRequest) (response.Response, error)
	LoginUser(request request.LoginRequest) (response.Response, error)
	GetUsers(request request.PaginationQuery) (response.Response, error)
	GetUser(id int) (response.Response, error)
	UpdateUser(id int, user request.UpdateUser) (response.Response, error)
	DeleteUser(id int) (response.Response, error)
	RefreshToken(refreshToken request.RefreshRequest) (response.Response, error)
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
	CreateProduct(request request.ProductRequest) response.Response
	GetProducts() response.Response
	GetProduct(id int) response.Response
	UpdateProduct(id int, product utils.UpdateProduct) response.Response
	DeleteProduct(id int) response.Response
}

type ProductItemService interface {
	CreateProductItem(request request.ProductItemRequest) response.Response
	GetProductItems() response.Response
	GetProductItem(id int) response.Response
	UpdateProductItem(id int, productItem utils.UpdateProductItem) response.Response
	DeleteProductItem(id int) response.Response
}

type ProductImageService interface {
	CreateProductImage(request request.ProductImageRequest) response.Response
}

type CartService interface {
	AddToCart(request request.CartRequest, userId uint) response.Response
}

type SizeService interface {
	CreateSize(request request.SizeRequest) response.Response
	GetSizes() response.Response
	GetSize(id int) response.Response
	UpdateSize(id int, size utils.UpdateSize) response.Response
	DeleteSize(id int) response.Response
}

type GetProductService interface {
	GetAllProducts() response.Response
	GetSingleProduct(id int) response.Response
}

type ReviewService interface {
	CreateReview(request request.ReviewRequest) response.Response
	GetReviews() response.Response
}
