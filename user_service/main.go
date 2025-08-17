package main

import (
	"main/controller"
	"main/database"
	"main/database/queries"
	"main/router"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	db := database.CreateConnection();
	uq := queries.UserQueryConstructor(db);
	controller := controller.NewController(uq);

	newRouter := router.NewRouter(app);

	newRouter.RegisterRouters(controller);

	app.Run(":5000");
}