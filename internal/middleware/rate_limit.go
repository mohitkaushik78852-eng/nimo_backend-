package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	started time.Time
	count   int
}

// RateLimit applies a fixed-window limit keyed by the connecting IP.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	entries := make(map[string]rateLimitEntry)

	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			ip = c.Request.RemoteAddr
		}

		now := time.Now()
		mu.Lock()
		for key, entry := range entries {
			if now.Sub(entry.started) >= window {
				delete(entries, key)
			}
		}
		entry := entries[ip]
		if entry.started.IsZero() || now.Sub(entry.started) >= window {
			entry = rateLimitEntry{started: now}
		}
		entry.count++
		entries[ip] = entry
		blocked := entry.count > limit
		mu.Unlock()

		if blocked {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many attempts, please try again later",
			})
			return
		}
		c.Next()
	}
}
