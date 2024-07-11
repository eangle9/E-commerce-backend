package jwttoken

import (
	"Eccomerce-website/internal/core/entity"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func GenerateTokenPair(id uint, role string) (map[string]string, error) {
	// viper.SetConfigFile("../.env")
	viper.AddConfigPath("../")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.ReadInConfig()
	secretKey := viper.Get("JWT_SECRET")
	if secretKey == nil {
		err := errors.New("unable to get jwt secret key")
		errorResponse := entity.AuthInternalError.Wrap(err, "failed to get jwt secret key in .env file").WithProperty(entity.StatusCode, 500)
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
		return nil, errorResponse
	}

	tokenMap := map[string]string{
		"access_token":  tokenString,
		"refresh_token": refreshTokenString,
	}

	return tokenMap, nil

}
