package middleware

import (
	lib "mew-gateway/internal/libs"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRateLimit(config *lib.Config) gin.HandlerFunc {
	if config == nil {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}

	rlCfg := config.Server.RateLimits
	if rlCfg.MaxRequests <= 0 ||
		!strings.EqualFold(rlCfg.Mode, "on") {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}

	window := rlCfg.UpdateIn
	if window <= 0 {
		window = time.Minute
	}

	type bucket struct {
		remaining int
		reset     time.Time
	}

	var (
		mu      sync.Mutex
		buckets = make(map[string]*bucket)
	)

	return func(ctx *gin.Context) {
		now := time.Now()
		clientKey := ctx.ClientIP()

		mu.Lock()
		b, ok := buckets[clientKey]
		if !ok || now.After(b.reset) {
			b = &bucket{
				remaining: rlCfg.MaxRequests,
				reset:     now.Add(window),
			}
			buckets[clientKey] = b
		}

		if b.remaining <= 0 {
			resetUnix := b.reset.Unix()
			retryAfter := int(time.Until(b.reset).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}
			mu.Unlock()

			ctx.Header("Retry-After", strconv.Itoa(retryAfter))
			ctx.Header("X-RateLimit-Limit", strconv.Itoa(rlCfg.MaxRequests))
			ctx.Header("X-RateLimit-Remaining", "0")
			ctx.Header("X-RateLimit-Reset", strconv.FormatInt(resetUnix, 10))

			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry_after": retryAfter,
			})
			return
		}

		b.remaining--
		remaining := b.remaining
		reset := b.reset
		mu.Unlock()

		ctx.Header("X-RateLimit-Limit", strconv.Itoa(rlCfg.MaxRequests))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		ctx.Header("X-RateLimit-Reset", strconv.FormatInt(reset.Unix(), 10))

		ctx.Next()

		// Opportunistic cleanup of expired buckets to avoid unbounded growth.
		var cleanup bool
		mu.Lock()
		if len(buckets) > 1024 {
			cleanup = true
		}
		mu.Unlock()

		if cleanup {
			now = time.Now()
			mu.Lock()
			for key, item := range buckets {
				if now.After(item.reset) {
					delete(buckets, key)
				}
			}
			mu.Unlock()
		}
	}
}
