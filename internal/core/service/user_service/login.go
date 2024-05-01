package service

import (
	jwttoken "Eccomerce-website/internal/core/common/jwt_token"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
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

func (u userService) LoginUser(request request.LoginRequest) response.Response {
	email := request.Email
	password := request.Password

	user, err := u.userRepo.Authentication(email, password)
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusUnauthorized,
			ErrorType:    errorcode.Unauthorized,
			ErrorMessage: err.Error(),
		}
		return errorResponse
	}

	tokenMap, err := jwttoken.GenerateTokenPair(uint(user.ID))
	if err != nil {
		errorResponse := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: "failed to generate token pair",
		}
		return errorResponse
	}

	tokenData := loginData{
		Username: user.Username,
		Email:    user.Email,
		TokenPair: tokenData{
			AccessToken:  tokenMap["access_token"],
			RefreshToken: tokenMap["refresh_token"],
		},
	}

	errorMessage := fmt.Sprintf("Login successful! Wellcome back, %s!", user.Username)
	response := response.Response{
		Data:         tokenData,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: errorMessage,
	}
	return response
}
