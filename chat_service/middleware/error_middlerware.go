package middleware

import (
	. "main/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware[T any](endpoint Endpoint[T]) gin.HandlerFunc {
	return func(context *gin.Context) {	
		if result, err := endpoint(context); err != nil {
			context.JSON(http.StatusInternalServerError, HttpResponse{
				Done: false,	
				Error: err.Error(),
			})
			context.Abort()
			return
		} else {
			context.JSON(http.StatusOK, HttpResponse{
				Done: true,
				Result: result,
			})
		} 
	}
}