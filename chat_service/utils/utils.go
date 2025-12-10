package utils

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

func Contains(slice pq.Int64Array, item int64) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

type Claims struct {
	UserID int    `json:"id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func DecodeJWT(jwt_token string) (*Claims, error) {
	if err := godotenv.Load(); err != nil {
		println("Error loading .env file:", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(
		jwt_token,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		println("Ошибка парсинга токена:", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Ошибка парсинга токена доступа")
	}
}

func ExtractClaimsFromCookie(context *gin.Context) (*Claims, error) {
	cookie := context.Request.Cookies()

	jwt_token := ""

	for _, val := range cookie {
		if val.Name == "session_token" {
			jwt_token = val.Value
		}
	}

	claims, err := DecodeJWT(jwt_token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

