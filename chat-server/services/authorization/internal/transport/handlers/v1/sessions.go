package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gox7/notify/services/authorization/internal/lib"
	"github.com/gox7/notify/services/authorization/models"
	"github.com/gox7/notify/services/authorization/pkg/tokens"
)

func HandleRefresh(config *models.LocalConfig, authService *lib.AuthorizathionService, sessionService *lib.SessionsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.TokenRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "invalid request body", nil,
			))
			return
		}

		if req.Token == "" {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "refresh token is required", nil,
			))
			return
		}

		// Validate refresh token (session)
		session, err := sessionService.SearchSession(req.Token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, ResponseJSON(
				"error", "invalid or expired refresh token", nil,
			))
			return
		}

		// Get user data
		user, err := authService.SearchUserByID(session.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseJSON(
				"error", "failed to retrieve user data", nil,
			))
			return
		}

		// Generate new access token
		accessToken, err := tokens.GenerateAccess(config, user.ID, user.Login, user.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseJSON(
				"error", "failed to generate access token", nil,
			))
			return
		}

		// Generate new refresh token
		newRefreshToken := tokens.GenerateRefresh()

		// Remove old session
		err = sessionService.RemoveSession(req.Token)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseJSON(
				"error", "failed to remove old session", nil,
			))
			return
		}

		// Create new session (refresh token expires in 7 days)
		expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()
		_, err = sessionService.CreateSession(user.ID, newRefreshToken, ctx.ClientIP(), expiresAt)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseJSON(
				"error", "failed to create new session", nil,
			))
			return
		}

		// Return new tokens
		ctx.JSON(http.StatusOK, ResponseJSON(
			"success", "tokens refreshed successfully",
			models.TokenData{
				AcessToken:   accessToken,
				RefreshToken: newRefreshToken,
			},
		))
	}
}

func HandleLogout(service *lib.SessionsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.TokenRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "invalid request body", nil,
			))
			return
		}

		if req.Token == "" {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "refresh token is required", nil,
			))
			return
		}

		// Remove session
		err := service.RemoveSession(req.Token)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseJSON(
				"error", "failed to remove session", nil,
			))
			return
		}

		ctx.JSON(http.StatusOK, ResponseJSON(
			"success", "logged out successfully", nil,
		))
	}
}

func HandleValidate(config *models.LocalConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.TokenRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "invalid request body", nil,
			))
			return
		}

		if req.Token == "" {
			ctx.JSON(http.StatusBadRequest, ResponseJSON(
				"error", "access token is required", nil,
			))
			return
		}

		// Validate access token
		claims, err := tokens.CheckAccess(config, req.Token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, ResponseJSON(
				"error", "invalid or expired access token", nil,
			))
			return
		}

		// Return user data
		ctx.JSON(http.StatusOK, ResponseJSON(
			"success", "token is valid",
			models.UserData{
				UserID: claims.UserId,
				Login:  claims.Login,
				Email:  claims.Email,
			},
		))
	}
}
