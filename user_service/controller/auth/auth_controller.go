package auth_controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"main/config"
	. "main/controller/dto"
	"main/controller/utils"
	"main/database/queries"
	"main/redis"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AQ  *queries.AuthQuery
	UQ  *queries.UserQuery
	RDB *redis.RedisConnector
}

func (ac *AuthController) LoginUser(context *gin.Context) {
	var userDTO UserDTO

	if err := context.ShouldBindJSON(&userDTO); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.AQ.GetUserByNameOrEmail(userDTO.Name, userDTO.Email)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   "GetUserByIDException",
			"message": err.Error(),
		})
		return
	}

	if user == nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Wrong Credentials",
		})
		return
	}

	if !utils.CheckPasswordHash(userDTO.Password, user.Password) {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Wrong Password",
		})
		return
	}

	jwtToken, err := utils.GenerateJWT(int64(user.ID), userDTO.Name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "CreatingScretKeyException",
			"message": err.Error(),
		})
	}

	context.SetCookie(
		"session_token",
		jwtToken,
		3600,
		"/",
		"localhost",
		false,
		true,
	)

	context.JSON(http.StatusOK, gin.H{
		"done": true,
	})
}

// func (ac *AuthController) RegisterUser(context *gin.Context) {
// 	var userDTO UserDTO

// 	if err := context.ShouldBindJSON(&userDTO); err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if userDTO.Name == "" || userDTO.Password == "" || userDTO.Email == "" {
// 		context.JSON(http.StatusUnauthorized, gin.H{
// 			"message": "Wrong Credentials",
// 		})
// 		return
// 	}

// 	hasedPassword, err := utils.HashPassword(userDTO.Password)

// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{
// 			"error":   "HashingError",
// 			"message": err.Error(),
// 		})
// 	}

// 	id, err := ac.UQ.InsertUser(userDTO.Email, userDTO.Name, hasedPassword)

// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{
// 			"error":   "createUserException",
// 			"message": err.Error(),
// 		})
// 	}

// 	jwtToken, err := utils.GenerateJWT(id, userDTO.Name)

// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{
// 			"error":   "CreatingScretKeyException",
// 			"message": err.Error(),
// 		})
// 	}

// 	context.SetCookie(
// 		"session_token",
// 		jwtToken,
// 		3600,
// 		"/",
// 		"localhost",
// 		false,
// 		true,
// 	)

// 	context.JSON(http.StatusOK, gin.H{
// 		"done": true,
// 	})
// }

func (ac *AuthController) SMTPApprove(context *gin.Context) {
	var userDTO UserDTO
	if err := context.ShouldBindJSON(&userDTO); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user, err := ac.AQ.GetUserByNameOrEmail(userDTO.Name, userDTO.Email); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "SMTPApprove",
			"message": err.Error(),
		})
		return
	} else if user != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "User with this credentials already exists",
		})
		return
	}

	if userDTO.Name == "" || userDTO.Password == "" || userDTO.Email == "" {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Wrong Credentials",
		})
		return
	}

	configFromContext, exists := context.Get("app_config")
	if !exists {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Config not found"})
		return
	}

	appConfig := configFromContext.(*config.Config)

	token := base64.StdEncoding.EncodeToString([]byte(userDTO.Email))
	verfifyLink := fmt.Sprintf("Verify Link: http://localhost:3000/verify/%s", token)

	if err := ac.RDB.SetData(token, userDTO); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "SMTPApprove",
			"message": err.Error(),
		})
		return
	}

	addr := fmt.Sprintf("%s:%s", appConfig.SMTP_ADDR, appConfig.SMTP_PORT)
	auth := smtp.PlainAuth("", appConfig.EMAIL_SENDER, appConfig.EMAIL_PASSWORD, appConfig.SMTP_ADDR)

	if err := smtp.SendMail(addr, auth, appConfig.EMAIL_SENDER, []string{userDTO.Email}, []byte(verfifyLink)); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "SMTPApprove",
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"done": true,
	})
}

func (ac *AuthController) VerifyResult(context *gin.Context) {
	var token TokenUriDto
	if err := context.ShouldBindUri(&token); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "VerifyResult",
			"message": err.Error(),
		})
		return
	}

	value, err := ac.RDB.GetData(token.Token)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "VerifyResult",
			"message": err.Error(),
		})
		return
	}

	var userDTO UserDTO

	if err := json.Unmarshal([]byte(value), &userDTO); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "VerifyResult",
			"message": err.Error(),
		})
		return
	}

	hasedPassword, err := utils.HashPassword(userDTO.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "HashingError",
			"message": err.Error(),
		})
	}

	id, err := ac.UQ.InsertUser(userDTO.Email, userDTO.Name, hasedPassword)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "createUserException",
			"message": err.Error(),
		})
		return
	}

	jwtToken, err := utils.GenerateJWT(id, userDTO.Name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "CreatingScretKeyException",
			"message": err.Error(),
		})
		return
	}

	context.SetCookie(
		"session_token",
		jwtToken,
		3600,
		"/",
		"localhost",
		false,
		true,
	)

	context.JSON(http.StatusOK, gin.H{
		"done": true,
	})
}

func (a *AuthController) Logout(context *gin.Context) {
	context.
		SetCookie(
			"session_token",
			"",
			0,
			"/",
			"localhost",
			false,
			true,
		)

	context.JSON(http.StatusOK, gin.H{
		"done": true,
	})
}
