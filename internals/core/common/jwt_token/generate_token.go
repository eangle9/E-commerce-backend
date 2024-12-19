package jwttoken

import (
	"Eccomerce-website/internals/core/entity"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GenerateTokenPair(ctx context.Context, id uint, role string, serviceLogger *zap.Logger, requestID string, client *redis.Client) (map[string]string, error) {
	// viper.SetConfigFile("../.env")
	viper.AddConfigPath("../")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.ReadInConfig()
	secretKey := viper.Get("JWT_SECRET")
	if secretKey == nil {
		err := errors.New("unable to get jwt secret key")
		errorResponse := entity.AuthInternalError.Wrap(err, "failed to get jwt secret key in .env file").WithProperty(entity.StatusCode, 500)
		serviceLogger.Error("failed to secret key in .env file",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "GenerateTokenPair"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}
	secret := secretKey.(string)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   id,
			"role": role,
			"exp":  time.Now().Add(15 * time.Minute).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		errorResponse := entity.AuthInternalError.Wrap(err, "failed to signed access_token by secret key").WithProperty(entity.StatusCode, 500)
		serviceLogger.Error("failed to get signed token",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "GenerateTokenPair"),
			zap.String("requestID", requestID),
			zap.Uint("userID", id),
			zap.String("role", role),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   id,
			"role": role,
			"exp":  time.Now().Add(30 * 24 * time.Hour).Unix(),
		})

	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		errorResponse := entity.AuthInternalError.Wrap(err, "failed to signed refresh_token by secret key").WithProperty(entity.StatusCode, 500)
		serviceLogger.Error("failed to get signed refreshToken",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "GenerateTokenPair"),
			zap.String("requestID", requestID),
			zap.Uint("userID", id),
			zap.String("role", role),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	refreshTokenExpiration := 30 * 24 * time.Hour
	if err = client.Set(ctx, refreshTokenString, "false", refreshTokenExpiration).Err(); err != nil {
		errorResponse := entity.AppInternalError.Wrap(err, "failed to store refresh token in redis").WithProperty(entity.StatusCode, 500)
		serviceLogger.Error("unable to store refresh token in redis",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "GenerateTokenPair"),
			zap.String("requestID", requestID),
			zap.String("refreshToken", refreshTokenString),
			zap.Uint("userID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	tokenMap := map[string]string{
		"access_token":  tokenString,
		"refresh_token": refreshTokenString,
	}

	return tokenMap, nil

}
