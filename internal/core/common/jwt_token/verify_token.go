package jwttoken

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func VerifyToken(tokenString string) (uint, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			return nil, err
		}
		viper.AddConfigPath("../")
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.ReadInConfig()
		secretKey := viper.Get("JWT_SECRET").(string)

		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		err := errors.New("invalid token")
		return 0, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err := errors.New("unable to extract claims from the token")
		return 0, "", err
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		err := errors.New("unable to get id from the claims")
		return 0, "", err
	}

	role, ok := claims["role"].(string)
	if !ok {
		err := errors.New("unable to get role from claims")
		return 0, "", err
	}

	id := uint(idFloat)

	return id, role, nil
}
