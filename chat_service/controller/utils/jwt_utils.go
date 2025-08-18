package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID int    `json:"id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func DecodeJWT(jwt_token string) (*Claims, error) {
	if err := godotenv.Load(); err != nil {
		println("Error loading .env file:", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET");

	token, err := jwt.ParseWithClaims(
		jwt_token,
		&Claims{}, 
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil 
		},
	)

	if err != nil {
		println("Ошибка парсинга токена:", err)
		return nil, err;
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil;
	} else {
		return nil, fmt.Errorf("Ошибка парсинга токена доступа")
	}
}