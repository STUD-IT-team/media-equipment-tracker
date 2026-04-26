package middleware

import (
	"time"

	"github.com/containerd/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		entry := logrus.WithFields(log.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"status":     c.Writer.Status(),
			"duration":   time.Since(startTime),
		})

		if c.Writer.Status() >= 500 {
			entry.Error(c.Errors.String())
		} else if c.Writer.Status() >= 400 {
			entry.Warn("Request completed with client error")
		} else {
			entry.Info("Request completed successfully")
		}
	}
}
