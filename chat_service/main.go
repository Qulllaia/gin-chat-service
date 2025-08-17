package main

import (
	"main/controller"
	"main/router"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default();

	controllerChat := controller.NewController();
	routerChat := router.NewRouter(app);
	routerChat.RegisterRouters(controllerChat);

	app.Run(":5050");
}