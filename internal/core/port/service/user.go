package service

import (
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
)

type UserService interface {
	SignUp(request request.SignUpRequest) response.Response
	LoginUser(request request.LoginRequest) response.Response
	GetUsers() response.Response
}
