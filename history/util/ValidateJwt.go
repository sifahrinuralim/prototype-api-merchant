package util

import (
	"errors"
	"history/config"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateAndExtractClaims(tokenString string) (jwt.MapClaims, error) {
	var jwtSecret = []byte(config.SecretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("token parse error: " + err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
