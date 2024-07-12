package service

import (
	jwttoken "Eccomerce-website/internal/core/common/jwt_token"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
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

func (u userService) LoginUser(ctx context.Context, request request.LoginRequest, requestID string) (response.Response, error) {
	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed login validation").WithProperty(entity.StatusCode, 400)
		u.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "LoginUser"),
			zap.String("requestID", requestID),
			zap.String("email", request.Email),
			zap.String("password_length", fmt.Sprintf("%d", len(request.Password))),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	user, err := u.userRepo.Authentication(ctx, request, requestID)
	if err != nil {
		return response.Response{}, err
	}

	tokenMap, err := jwttoken.GenerateTokenPair(ctx, uint(user.ID), user.Role, u.serviceLogger, requestID)
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
