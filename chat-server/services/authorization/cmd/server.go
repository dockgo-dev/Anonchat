package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gox7/notify/services/authorization/internal/lib"
	"github.com/gox7/notify/services/authorization/internal/repository"
	"github.com/gox7/notify/services/authorization/internal/transport"
	"github.com/gox7/notify/services/authorization/models"
)

func main() {
	config := new(models.LocalConfig)
	lib.NewConfig(config)

	postgreslog := new(slog.Logger)
	lib.NewLogger("server", postgreslog)

	postgres := new(repository.Postgres)
	repository.NewPostgres(config, postgreslog, postgres)
	postgres.Migration()

	authService := new(lib.AuthorizathionService)
	sessionService := new(lib.SessionsService)
	lib.NewAuthorizathion(postgres, authService)
	lib.NewSessions(postgres, sessionService)

	engine := gin.Default()
	transport.Register(config, &transport.Services{
		Auth:    authService,
		Session: sessionService,
	}, engine)
	transport.Listen(config, engine)
}
