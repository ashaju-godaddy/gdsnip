package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string]int
	limit    int
	mu       sync.RWMutex
}

// NewRateLimiter creates a new rate limiter with the specified requests per minute
func NewRateLimiter(limit int) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
	}

	// Start cleanup goroutine that resets counters every minute
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.mu.Lock()
			rl.requests = make(map[string]int)
			rl.mu.Unlock()
		}
	}()

	return rl
}

// Middleware returns Echo middleware function
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get client IP
			ip := c.RealIP()

			// Check current request count
			rl.mu.Lock()
			count := rl.requests[ip]
			if count >= rl.limit {
				rl.mu.Unlock()
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"success": false,
					"error": map[string]interface{}{
						"code":    "RATE_LIMIT_EXCEEDED",
						"message": "too many requests, please try again later",
					},
				})
			}
			rl.requests[ip]++
			rl.mu.Unlock()

			return next(c)
		}
	}
}
