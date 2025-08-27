package utils

import "github.com/gin-gonic/gin"

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