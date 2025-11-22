package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gox7/notify/services/authorization/internal/lib"
	"github.com/gox7/notify/services/authorization/models"
	"github.com/gox7/notify/services/authorization/pkg/tokens"
)

func HandleRegister(config *models.LocalConfig, authService *lib.AuthorizathionService, sessionService *lib.SessionsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.RegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "invalid request body", nil,
			))
			return
		}

		// Create user
		userID, err := authService.CreateUser(req.Login, req.Email, req.Password, ctx.ClientIP())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", err.Error(), nil,
			))
			return
		}

		// Get user data to generate tokens
		user, err := authService.SearchUserByID(userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseJSON(
				"error", "failed to retrieve user data", nil,
			))
			return
		}

		// Generate tokens
		accessToken, err := tokens.GenerateAccess(config, user.ID, user.Login, user.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseJSON(
				"error", "failed to generate access token", nil,
			))
			return
		}

		refreshToken := tokens.GenerateRefresh()

		// Create session (refresh token expires in 7 days)
		expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()
		_, err = sessionService.CreateSession(user.ID, refreshToken, ctx.ClientIP(), expiresAt)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseJSON(
				"error", "failed to create session", nil,
			))
			return
		}

		// Return tokens
		ctx.JSON(http.StatusCreated, ResponseJSON(
			"success", "user registered successfully",
			models.TokenData{
				AcessToken:   accessToken,
				RefreshToken: refreshToken,
			},
		))
	}
}

func HandleLogin(config *models.LocalConfig, authService *lib.AuthorizathionService, sessionService *lib.SessionsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "invalid request body", nil,
			))
			return
		}

		// Search and validate user
		user, err := authService.SearchUser(req.Email, req.Password)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, ResponseJSON(
				"error", "invalid login or password", nil,
			))
			return
		}

		// Generate tokens
		accessToken, err := tokens.GenerateAccess(config, user.ID, user.Login, user.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseJSON(
				"error", "failed to generate access token", nil,
			))
			return
		}

		refreshToken := tokens.GenerateRefresh()

		// Create session (refresh token expires in 7 days)
		expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()
		_, err = sessionService.CreateSession(user.ID, refreshToken, ctx.ClientIP(), expiresAt)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseJSON(
				"error", "failed to create session", nil,
			))
			return
		}

		// Return tokens
		ctx.JSON(http.StatusOK, ResponseJSON(
			"success", "login successful",
			models.TokenData{
				AcessToken:   accessToken,
				RefreshToken: refreshToken,
			},
		))
	}
}
