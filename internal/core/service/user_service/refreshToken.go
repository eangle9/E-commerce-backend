package service

import (
	jwttoken "Eccomerce-website/internal/core/common/jwt_token"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

type tokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u userService) RefreshToken(refreshToken request.RefreshRequest) (response.Response, error) {
	if err := refreshToken.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed validation").WithProperty(entity.StatusCode, 400)
		return response.Response{}, errorResponse
	}

	rfToken := refreshToken.RefreshToken
	id, role, err := jwttoken.VerifyToken(rfToken)
	if err != nil {
		return response.Response{}, err
	}

	tokenMap, err := jwttoken.GenerateTokenPair(id, role)
	if err != nil {
		return response.Response{}, err
	}

	tokenPair := tokenPair{
		AccessToken:  tokenMap["access_token"],
		RefreshToken: tokenMap["refresh_token"],
	}

	response := response.Response{
		Data:       tokenPair,
		StatusCode: http.StatusOK,
		Message:    "Token Refresh Successful: Your access and refresh tokens have been successfully refreshed",
	}
	return response, nil
}
