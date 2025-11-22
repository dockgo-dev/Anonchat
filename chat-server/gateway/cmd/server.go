package main

import (
	"fmt"
	lib "mew-gateway/internal/libs"
	"mew-gateway/internal/transport"
	"mew-gateway/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := lib.LoadConfig()
	if err != nil {
		fmt.Println("aw, config loading:", err.Error())
	}

	logger, loggerCancel, err := lib.NewLogger("server")
	if err != nil {
		fmt.Println("aw, logger creating:", err.Error())
	}
	defer loggerCancel()

	hub := websocket.NewHub()
	go hub.Run()

	engine := gin.Default()
	transport.Register(engine, hub, logger, config)
	transport.Listen(engine, config)
}
