package auth_controller

import (
	. "main/controller/dto"
	"main/controller/utils"
	"main/database/queries"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AQ *queries.AuthQuery
	UQ *queries.UserQuery
};

func (ac *AuthController) LoginUser(context *gin.Context) {

	var userDTO UserDTO;

	if err := context.ShouldBindJSON(&userDTO); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.AQ.GetUserByName(userDTO.Name);

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "GetUserByIDException",
			"message": err.Error(),
		})
		return
	}

	if !utils.CheckPasswordHash(userDTO.Password, user.Password) {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Wrong Password",
		})
		return
	}

	jwtToken, err := utils.GenerateJWT(userDTO.ID, userDTO.Name);

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "CreatingScretKeyException",
			"message": err.Error(),
		})
	}
	
    context.SetCookie(
        "session_token",
        jwtToken,
        3600,
        "/",
        "",
        true,
        true,
    )

	context.JSON(http.StatusOK, gin.H{
		"done": true,
	})
}

func (ac *AuthController) RegisterUser(context *gin.Context) {
	var userDTO UserDTO;

	if err := context.ShouldBindJSON(&userDTO); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userDTO.Name == "" || userDTO.Password == "" {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Wrong Credentials",
		})
		return
	}

	hasedPassword, err := utils.HashPassword(userDTO.Password);

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "HashingError",
			"message": err.Error(),
		})
	}

	err = ac.UQ.InsertUser(userDTO.Name, hasedPassword);

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "createUserException",
			"message": err.Error(),
		})
	}

	jwtToken, err := utils.GenerateJWT(userDTO.ID, userDTO.Name);

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "CreatingScretKeyException",
			"message": err.Error(),
		})
	}

    context.SetCookie(
        "session_token",
        jwtToken,
        3600,
        "/",
        "",
        true,
        true,
    )

	context.JSON(http.StatusOK, gin.H{
		"done": true,
	})
}