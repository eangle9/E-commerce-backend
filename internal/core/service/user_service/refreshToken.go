package service

import (
	jwttoken "Eccomerce-website/internal/core/common/jwt_token"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

type tokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u userService) RefreshToken(refreshToken request.RefreshRequest) response.Response {
	rfToken := refreshToken.RefreshToken

	id, err := jwttoken.VerifyToken(rfToken)
	if err != nil {
		response := response.Response{
			Status:       http.StatusUnauthorized,
			ErrorType:    errorcode.Unauthorized,
			ErrorMessage: err.Error(),
		}
		return response
	}

	tokenMap, err := jwttoken.GenerateTokenPair(id)
	if err != nil {
		response := response.Response{
			Status:       http.StatusInternalServerError,
			ErrorType:    errorcode.InternalError,
			ErrorMessage: err.Error(),
		}
		return response
	}

	tokenPair := tokenPair{
		AccessToken:  tokenMap["access_token"],
		RefreshToken: tokenMap["refresh_token"],
	}

	response := response.Response{
		Data:         tokenPair,
		Status:       http.StatusOK,
		ErrorType:    errorcode.Success,
		ErrorMessage: "Token Refresh Successful: Your access and refresh tokens have been successfully refreshed",
	}
	return response
}
