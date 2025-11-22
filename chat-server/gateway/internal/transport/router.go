package transport

import (
	"fmt"
	"log/slog"
	lib "mew-gateway/internal/libs"
	v1 "mew-gateway/internal/transport/handlers/v1"
	"mew-gateway/internal/transport/middleware"
	"mew-gateway/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine, hub *websocket.Hub, logger *slog.Logger, config *lib.Config) {
	if config.Server.Logger.Mode == "on" {
		engine.Use(
			middleware.NewLogger(logger, config),
		)
	}
	if config.Server.RateLimits.Mode == "on" {
		engine.Use(
			middleware.NewRateLimit(config),
		)
	}
	if config.Server.Auth.Mode == "on" {
		engine.Use(
			middleware.NewAuth(config),
		)
	}

	engine.GET("/api/v1/status", v1.HandleStatus)
	auth := engine.Group("/api/v1/auth")
	{
		auth.POST("/register", v1.HandleAuthRegister(config))
		auth.POST("/login", v1.HandleAuthLogin(config))
		auth.POST("/refresh", v1.HandleAuthRefresh(config))
		auth.POST("/logout", v1.HandleAuthLogout(config))
		auth.GET("/validate", v1.HandleAuthValidate(config))
	}

	rooms := engine.Group("/api/v1/rooms")
	{
		rooms.GET("/ws", v1.HandleWebsocketRoom(config, hub))
		rooms.POST("/create", v1.HandleCreateRoom(hub))
		rooms.DELETE("/remove", v1.HandleRemoveRoom(hub))
	}
}

func Listen(engine *gin.Engine, config *lib.Config) {
	server := http.Server{
		Handler:      engine,
		Addr:         config.Server.Addr,
		WriteTimeout: config.Server.Timeouts.Write,
		ReadTimeout:  config.Server.Timeouts.Read,
	}

	fmt.Println("[+] server.listen:", config.Server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("[-] server.listen:", err.Error())
	}
}
