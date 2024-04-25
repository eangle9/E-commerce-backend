package jwttoken

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func GenerateTokenPair(id uint) (map[string]string, error) {
	// viper.SetConfigFile("../.env")
	viper.AddConfigPath("../")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.ReadInConfig()
	secretKey := viper.Get("JWT_SECRET")
	if secretKey == nil {
		err := errors.New("unable to get jwt secret key")
		return nil, err
	}
	secret := secretKey.(string)
	fmt.Println("secretKey: ", secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  id,
			"exp": time.Now().Add(30 * time.Second).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  id,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		})

	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	tokenMap := map[string]string{
		"access_token":  tokenString,
		"refresh_token": refreshTokenString,
	}

	return tokenMap, nil

}
