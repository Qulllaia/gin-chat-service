package middleware

import (
	"main/config"

	"github.com/gin-gonic/gin"
)

func ConfigMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("app_config", config)
		c.Next()
	}
}
