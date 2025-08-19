package main

import (
	"main/controller"
	"main/database"
	"main/router"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default();

	db := database.CreateConnection();

	controllerChat := controller.NewController(db);
	routerChat := router.NewRouter(app);
	routerChat.RegisterRouters(controllerChat);

	app.Run(":5050");
}