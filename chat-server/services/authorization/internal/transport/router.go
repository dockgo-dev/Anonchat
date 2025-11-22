package transport

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gox7/notify/services/authorization/internal/lib"
	v1 "github.com/gox7/notify/services/authorization/internal/transport/handlers/v1"
	"github.com/gox7/notify/services/authorization/models"
)

type (
	Services struct {
		Auth    *lib.AuthorizathionService
		Session *lib.SessionsService
	}
)

func Register(config *models.LocalConfig, s *Services, engine *gin.Engine) {
	engine.GET("/v1/status", v1.HandleStatus)

	engine.POST("/v1/register", v1.HandleRegister(config, s.Auth, s.Session))
	engine.POST("/v1/login", v1.HandleLogin(config, s.Auth, s.Session))
	engine.POST("/v1/refresh", v1.HandleRefresh(config, s.Auth, s.Session))
	engine.POST("/v1/logout", v1.HandleLogout(s.Session))
	engine.POST("/v1/validate", v1.HandleValidate(config))
}

func Listen(config *models.LocalConfig, engine *gin.Engine) {
	server := http.Server{
		Handler:      engine,
		Addr:         config.Server.Addr,
		WriteTimeout: config.Server.Timeouts.Write,
		ReadTimeout:  config.Server.Timeouts.Read,
	}

	fmt.Println("[+] transport.listen:", config.Server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("[-] transport.listen:", err)
	}
}
