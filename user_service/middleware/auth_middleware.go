package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func AuthMiddleware(context *gin.Context) {
	jwtCookie, err := context.Cookie("session_token");
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "session_token error",
			"message": err.Error(),
		})
		context.Abort()
		return
	}

	if jwtCookie == "" {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "No session_token",
		})
		context.Abort()
		return
	}

	if err := godotenv.Load(); err != nil {
        println("Error loading .env file:", err)
    }

	jwtSecret := os.Getenv("JWT_SECRET");

	token, err := jwt.Parse(jwtCookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	});

	if err != nil || !token.Valid {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Token",
		})
		context.Abort()
		return
	}
}