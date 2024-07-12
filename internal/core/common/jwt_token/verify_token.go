package jwttoken

import (
	"Eccomerce-website/internal/core/entity"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func VerifyToken(ctx context.Context, tokenString string, serviceLogger *zap.Logger, requestID string) (uint, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			errorResponse := entity.AppInternalError.Wrap(err, "incorrect signing method").WithProperty(entity.StatusCode, 500)
			serviceLogger.Error("unexpected signing method",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "serviceLayer"),
				zap.String("function", "VerifyToken"),
				zap.String("requestID", requestID),
				zap.Any("method", t.Header["alg"]),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return nil, errorResponse
		}
		viper.AddConfigPath("../")
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.ReadInConfig()
		secretKey := viper.Get("JWT_SECRET").(string)

		return []byte(secretKey), nil
	})

	if err != nil {
		errorResponse := entity.InvalidCredentials.Wrap(err, "invalid token").WithProperty(entity.StatusCode, 401)
		serviceLogger.Error("invalid token",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "VerifyToken"),
			zap.String("requestID", requestID),
			zap.String("tokenString", tokenString),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return 0, "", errorResponse
	}

	if !token.Valid {
		err := errors.New("invalid or expired token")
		errorResponse := entity.InvalidCredentials.Wrap(err, "incorrect token").WithProperty(entity.StatusCode, 401)
		serviceLogger.Error("invalid or expired token",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "VerifyToken"),
			zap.String("requestID", requestID),
			zap.String("tokenString", tokenString),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return 0, "", errorResponse
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err := errors.New("unable to extract claims from the token")
		errorResponse := entity.AuthInternalError.Wrap(err, "failed to get token claims").WithProperty(entity.StatusCode, 500)
		serviceLogger.Error("failed to get token claims",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "VerifyToken"),
			zap.String("requestID", requestID),
			zap.String("tokenString", tokenString),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return 0, "", errorResponse
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		err := errors.New("unable to get id from the claims")
		errorResponse := entity.AuthInternalError.Wrap(err, "failed to get userID in the claims").WithProperty(entity.StatusCode, 500)
		serviceLogger.Error("failed to extract userID from the token claims",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "VerifyToken"),
			zap.String("requestID", requestID),
			zap.String("tokenString", tokenString),
			zap.Any("claims", claims),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return 0, "", errorResponse
	}

	role, ok := claims["role"].(string)
	if !ok {
		err := errors.New("unable to get role from claims")
		errorResponse := entity.AuthInternalError.Wrap(err, "failed to get role in the claims").WithProperty(entity.StatusCode, 500)
		serviceLogger.Error("failed to extract user_role from the token claims",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "VerifyToken"),
			zap.String("requestID", requestID),
			zap.String("tokenString", tokenString),
			zap.Any("claims", claims),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return 0, "", errorResponse
	}

	id := uint(idFloat)

	return id, role, nil
}
