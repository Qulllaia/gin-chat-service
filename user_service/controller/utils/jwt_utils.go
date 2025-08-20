package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID int64 `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int64, name string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	//println(userID)
	claims := &Claims{
		UserID: userID,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	if err := godotenv.Load(); err != nil {
        println("Error loading .env file:", err)
    }

	jwtSecret := os.Getenv("JWT_SECRET");

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}