package service

import (
	jwttoken "Eccomerce-website/internal/core/common/jwt_token"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

type loginData struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	TokenPair tokenData `json:"tokenPair"`
}

type tokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u userService) LoginUser(request request.LoginRequest) (response.Response, error) {
	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed login validation").WithProperty(entity.StatusCode, 400)
		return response.Response{}, errorResponse
	}

	user, err := u.userRepo.Authentication(request)
	if err != nil {
		return response.Response{}, err
	}

	tokenMap, err := jwttoken.GenerateTokenPair(uint(user.ID), user.Role)
	if err != nil {
		return response.Response{}, err
	}

	tokenData := loginData{
		Username: user.Username,
		Email:    user.Email,
		TokenPair: tokenData{
			AccessToken:  tokenMap["access_token"],
			RefreshToken: tokenMap["refresh_token"],
		},
	}

	response := response.Response{
		Data:       tokenData,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("Login successful! Wellcome back, %s!", user.Username),
	}
	return response, nil
}
