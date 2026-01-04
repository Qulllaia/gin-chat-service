package main

import (
	"time"

	"main/config"
	"main/controller"
	"main/database"
	"main/database/queries"
	"main/redis"
	"main/router"
	"main/user"

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

	config, err := config.CreateConfig()
	if err != nil {
		println(err.Error())
		panic("CONFIG ERROR")
	}
	redisConnection := redis.NewRedisConnector()
	defer redisConnection.Close()

	db, err := database.CreateConnection(config)
	if err != nil {
		panic(err.Error())
	}
	uq := queries.UserQueryConstructor(db)
	aq := queries.AuthQueryConstructor(db)
	controller := controller.NewController(uq, aq, redisConnection)

	go user.StartUserServer(uq)

	newRouter := router.NewRouter(app)

	newRouter.RegisterRouters(controller, config)
	app.Run(":5000")
}
