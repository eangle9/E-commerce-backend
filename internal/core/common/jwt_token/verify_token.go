package jwttoken

import (
	"Eccomerce-website/internal/core/entity"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func VerifyToken(tokenString string) (uint, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			errorResponse := entity.AppInternalError.Wrap(err, "incorrect signing method").WithProperty(entity.StatusCode, 500)
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
		return 0, "", errorResponse
	}

	if !token.Valid {
		err := errors.New("invalid or expired token")
		errorResponse := entity.InvalidCredentials.Wrap(err, "incorrect token").WithProperty(entity.StatusCode, 401)
		return 0, "", errorResponse
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err := errors.New("unable to extract claims from the token")
		errorResponse := entity.AuthInternalError.Wrap(err, "failed to get token claims").WithProperty(entity.StatusCode, 500)
		return 0, "", errorResponse
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		err := errors.New("unable to get id from the claims")
		errorResponse := entity.AuthInternalError.Wrap(err, "failed to get userID in the claims").WithProperty(entity.StatusCode, 500)
		return 0, "", errorResponse
	}

	role, ok := claims["role"].(string)
	if !ok {
		err := errors.New("unable to get role from claims")
		errorResponse := entity.AuthInternalError.Wrap(err, "failed to get role in the claims").WithProperty(entity.StatusCode, 500)
		return 0, "", errorResponse
	}

	id := uint(idFloat)

	return id, role, nil
}
