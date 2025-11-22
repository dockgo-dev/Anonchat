package middleware

import (
	"log/slog"
	lib "mew-gateway/internal/libs"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func NewLogger(logger *slog.Logger, config *lib.Config) gin.HandlerFunc {
	if logger == nil {
		logger = slog.Default()
	}

	var (
		modeEnabled = true
		baseLevel   = slog.LevelInfo
	)

	if config != nil {
		loggerCfg := config.Server.Logger
		modeEnabled = strings.EqualFold(loggerCfg.Mode, "on") || loggerCfg.Mode == ""
		if parsed, ok := parseLevel(loggerCfg.Level); ok {
			baseLevel = parsed
		}
	}

	if !modeEnabled {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}

	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path

		ctx.Next()

		status := ctx.Writer.Status()
		latency := time.Since(start)

		level := baseLevel
		switch {
		case status >= http.StatusInternalServerError:
			level = slog.LevelError
		case status >= http.StatusBadRequest && level < slog.LevelWarn:
			level = slog.LevelWarn
		}

		if level < baseLevel {
			level = baseLevel
		}

		if !logger.Enabled(ctx, level) {
			return
		}

		logger.LogAttrs(
			ctx,
			level,
			"http_request",
			slog.Int("status", status),
			slog.String("method", ctx.Request.Method),
			slog.String("path", path),
			slog.String("client_ip", ctx.ClientIP()),
			slog.String("user_agent", ctx.Request.UserAgent()),
			slog.Duration("latency", latency),
			slog.Int64("elapsed_ms", latency.Milliseconds()),
		)
	}
}

func parseLevel(value string) (slog.Level, bool) {
	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug, true
	case "info":
		return slog.LevelInfo, true
	case "warn", "warning":
		return slog.LevelWarn, true
	case "err", "error":
		return slog.LevelError, true
	default:
		return slog.LevelInfo, false
	}
}
