package service

import (
	jwttoken "Eccomerce-website/internal/core/common/jwt_token"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type tokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u userService) RefreshToken(ctx context.Context, refreshToken request.RefreshRequest, requestID string) (response.Response, error) {
	if err := refreshToken.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed validation").WithProperty(entity.StatusCode, 400)
		u.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "RefreshToken"),
			zap.String("requestID", requestID),
			zap.Any("requestData", refreshToken),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	rfToken := refreshToken.RefreshToken
	id, role, err := jwttoken.VerifyRefreshToken(ctx, rfToken, u.serviceLogger, requestID, u.client)
	if err != nil {
		return response.Response{}, err
	}

	tokenMap, err := jwttoken.GenerateTokenPair(ctx, id, role, u.serviceLogger, requestID, u.client)
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
