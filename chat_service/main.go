package main

import (
	"main/websockets"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default();

	app.GET("/ws", websockets.WebsocketsInit)

	app.Run(":5050");
}