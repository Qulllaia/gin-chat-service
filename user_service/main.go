package main

import (
	"main/controller"
	"main/router"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	controller := controller.NewController();

	newRouter := router.NewRouter(app);

	newRouter.RegisterRouters(controller);

	app.Run(":5000");
}