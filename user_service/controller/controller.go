package controller

import (
	auth_controller "main/controller/auth"
	user_controller "main/controller/user"
)

type Controller struct {
	Auth auth_controller.AuthController;
	User user_controller.UserController;
}

func NewController() *Controller{
	return &Controller{
		Auth: auth_controller.AuthController{},
		User: user_controller.UserController{},
	};
}