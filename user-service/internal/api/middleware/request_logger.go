package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"log/slog"
)

func RequestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		log.Info("HTTP request",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.String("ip", clientIP),
			slog.Duration("duration", duration),
		)
	}
}
