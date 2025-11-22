package middleware

import (
	"mew-gateway/internal/app"
	lib "mew-gateway/internal/libs"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewAuth(config *lib.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		access := ctx.GetHeader("Authorization")
		if access != "" {
			token := access
			if strings.HasPrefix(access, "Bearer ") {
				token = strings.TrimPrefix(access, "Bearer ")
			}

			data, err := app.Validate(config, token)
			if err == nil {
				ctx.Set("user.id", data.UserID)
				ctx.Set("user.login", data.Login)
				ctx.Set("user.email", data.Email)
			}
		}
		ctx.Next()
	}
}
