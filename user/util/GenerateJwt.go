package util

import (
	"time"
	"user/config"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(credential, cif string) (string, error) {
	var jwtSecret = []byte(config.SecretKey)

	claims := jwt.MapClaims{
		"credential": credential,
		"cif":        cif,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
