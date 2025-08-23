package main

import (
	"main/controller"
	"main/database"
	"main/database/queries"
	"main/router"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default();

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Кеширование CORS-префлайта
	}))
	
	db := database.CreateConnection();

	cq := queries.ChatQueryConstructor(db);

	controllerChat := controller.NewController(db, cq);
	routerChat := router.NewRouter(app);
	routerChat.RegisterRouters(controllerChat);

	app.Run(":5050");
}
