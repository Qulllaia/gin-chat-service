package middleware

import (
	. "main/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware[T any](endpoint Endpoint[T]) gin.HandlerFunc {
	return func(context *gin.Context) {	
		if errType, result, err := endpoint(context); errType != NoError {
			context.JSON(http.StatusInternalServerError, HttpResponse{
				Done: false,
				ErrorType: errType,	
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