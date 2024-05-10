package service

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
)

type UserService interface {
	SignUp(request request.SignUpRequest) response.Response
	LoginUser(request request.LoginRequest) response.Response
	GetUsers() response.Response
	GetUser(id int) response.Response
	UpdateUser(id int, user utils.UpdateUser) response.Response
	DeleteUser(id int) response.Response
	RefreshToken(refreshToken request.RefreshRequest) response.Response
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
