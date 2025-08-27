package main

import (
	"context"
	"main/config"
	"main/controller"
	"main/database"
	"main/database/queries"
	"main/router"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Кеширование CORS-префлайта
	}))

	ctx := context.Background()
    
	config, err := config.CreateConfigDatabase();
    
	ctx = context.WithValue(ctx, "config", config)
	

	db, err := database.CreateConnection(ctx);
	uq := queries.UserQueryConstructor(db);
	aq := queries.AuthQueryConstructor(db);
	controller := controller.NewController(uq, aq);

	newRouter := router.NewRouter(app);

	newRouter.RegisterRouters(controller);
	if err == nil {
		app.Run(":5000");
	}
}