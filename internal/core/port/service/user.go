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

type ProductCategorySevice interface {
	CreateProductCategory(request request.ProductCategoryRequest) response.Response
}
