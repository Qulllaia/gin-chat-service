package user_controller

import (
	. "main/controller/dto"
	"main/controller/utils"
	"main/database/queries"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct{
	UQ *queries.UserQuery
};

func (uc *UserController) CreateUser(context *gin.Context) {

	var user UserDTO;

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} 

	_, err := uc.UQ.InsertUser(user.Name, user.Password);
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "createUserException",
			"message": err.Error(),
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"done": true,
		})
	}
}

func (uc *UserController) GetAllUsers(context *gin.Context) {
	users, err := uc.UQ.GetAllUsers();
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "GetAllUsersException",
			"message": err.Error(),
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"done": true,
			"result": users,
		})
	}
}

func (uc *UserController) GetUserByID(context *gin.Context) {

	var userID UserIDURI;

	if err := context.ShouldBindUri(&userID); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := uc.UQ.GetUserByID(userID.ID);
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "GetUserByIDException",
			"message": err.Error(),
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"done": true,
			"result": users,
		})
	}
}

func (uc *UserController) UpdateUser(context *gin.Context) {

	var user UserDTO;

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} 

	err := uc.UQ.UpdateUser(user.ID, user.Name, user.Password);
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "UpdateUserException",
			"message": err.Error(),
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"done": true,
		})
	}
}

func (uc *UserController) DeleteUser(context *gin.Context) {
	var userID UserIDURI;

	if err := context.ShouldBindUri(&userID); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.UQ.DeleteUser(userID.ID);
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "DeleteUserException",
			"message": err.Error(),
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"done": true,
		})
	}
}

func (uc *UserController) GetUsersFriends(context *gin.Context) {

	cookie := context.Request.Cookies();

	jwt_token := "";

	for _, val := range cookie {
		if val.Name == "session_token" {
			jwt_token = val.Value;
		}
	}

	claims, err := utils.DecodeJWT(jwt_token);

	users, err := uc.UQ.GetUsersFriends(int(claims.UserID));
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "GetUsersFriendsException",
			"message": err.Error(),
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"done": true,
			"result": users,
		})
	}
}