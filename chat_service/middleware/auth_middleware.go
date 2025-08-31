package middleware

import (
	"main/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	configFromContext, exists := context.Get("app_config")
	if !exists {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Config not found"})
		context.Abort()
		return
	}
	
	appConfig := configFromContext.(*config.Config)

	jwtSecret := appConfig.JWT_SECRET;

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