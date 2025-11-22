package v1

import (
	"mew-gateway/internal/app"
	lib "mew-gateway/internal/libs"
	"mew-gateway/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleAuthRegister(config *lib.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request RegisterRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "request body(JSON) is invalid", nil,
			))
			return
		}

		if err := validator.ValidateLogin(request.Login); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", err.Error(), nil,
			))
			return
		}

		if err := validator.ValidatePassword(request.Password); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", err.Error(), nil,
			))
			return
		}

		if err := validator.ValidateEmail(request.Email); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", err.Error(), nil,
			))
			return
		}

		accessToken, refreshToken, err := app.RegisterUser(config, request.Login, request.Email, request.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", err.Error(), nil,
			))
			return
		}

		ctx.JSON(http.StatusOK, ResponseJSON(
			"success", "user registered successfully",
			gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		))
	}
}

func HandleAuthLogin(config *lib.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request LoginRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "request body(JSON) is invalid", nil,
			))
			return
		}

		if err := validator.ValidateEmail(request.Email); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", err.Error(), nil,
			))
			return
		}
		if err := validator.ValidatePassword(request.Password); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", err.Error(), nil,
			))
			return
		}

		accessToken, refreshToken, err := app.LoginUser(config, request.Email, request.Password)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, ResponseJSON(
				"error", err.Error(), nil,
			))
			return
		}

		ctx.JSON(http.StatusOK, ResponseJSON(
			"success", "login successful",
			gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		))
	}
}

func HandleAuthRefresh(config *lib.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request TokenRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "request body(JSON) is invalid", nil,
			))
			return
		}

		if request.Token == "" {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "refresh token is required", nil,
			))
			return
		}

		accessToken, refreshToken, err := app.Refresh(config, request.Token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, ResponseJSON(
				"error", err.Error(), nil,
			))
			return
		}

		ctx.JSON(http.StatusOK, ResponseJSON(
			"success", "tokens refreshed successfully",
			gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		))
	}
}

func HandleAuthLogout(config *lib.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request TokenRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "request body(JSON) is invalid", nil,
			))
			return
		}

		if request.Token == "" {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "refresh token is required", nil,
			))
			return
		}

		err := app.Logout(config, request.Token)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseJSON(
				"error", err.Error(), nil,
			))
			return
		}

		ctx.JSON(http.StatusOK, ResponseJSON(
			"success", "logged out successfully", nil,
		))
	}
}

func HandleAuthValidate(config *lib.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetInt64("user.id")
		login := ctx.GetString("user.login")
		email := ctx.GetString("user.email")

		if id == 0 {
			ctx.JSON(http.StatusUnauthorized, ResponseJSON(
				"error", "none access token", nil,
			))
			return
		}

		ctx.JSON(http.StatusOK, ResponseJSON(
			"success", "token is valid", UserData{
				Id: id, Email: email,
				Login: login,
			},
		))
	}
}
