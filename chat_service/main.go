package main

import (
	"main/config"
	"main/controller"
	"main/database"
	"main/database/queries"
	"main/redis"
	"main/router"
	"main/user"
	"main/websockets"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default();

    app.Static("/background", "./static/background")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Кеширование CORS-префлайта
	}))

	config, err := config.CreateConfig();

	if err != nil {
		panic("CONFIG ERROR")
	}
	userServiceAddr := "localhost:50051"
	server, err := user.ConnectSerivce(userServiceAddr);
	
	if err != nil {
		println(err)
	}

	redisConnection := redis.NewRedisConnector();
	defer redisConnection.Close();
	
	db, _ := database.CreateConnection(config);
	defer db.DB.Close()

	
	cq := queries.ChatQueryConstructor(db, server);
	wsq := queries.WSQueryConstructor(db);
	
	connectorActor := websockets.NewConnectorActor(wsq, redisConnection)
    defer connectorActor.Stop()

	controllerChat := controller.NewController(db, cq, connectorActor, redisConnection);
	routerChat := router.NewRouter(app);
	routerChat.RegisterRouters(controllerChat, config);

	app.Run(":5050");
}
