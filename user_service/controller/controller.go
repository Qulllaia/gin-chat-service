package controller

import (
	auth_controller "main/controller/auth"
	user_controller "main/controller/user"
	"main/database/queries"
	"main/redis"
)

type Controller struct {
	Auth auth_controller.AuthController
	User user_controller.UserController
}

func NewController(uq *queries.UserQuery, aq *queries.AuthQuery, rdb *redis.RedisConnector) *Controller {

	return &Controller{
		Auth: auth_controller.AuthController{AQ: aq, UQ: uq, RDB: rdb},
		User: user_controller.UserController{UQ: uq},
	}
}
