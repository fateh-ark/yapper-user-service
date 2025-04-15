package controller

import (
	"fateh-ark/yapper-user-service/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerHandler(loggerInstance logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()

		logData := logger.LogData{
			Timestamp: &end,
			Level:     logger.InfoLogLevel,
			Component: "controller",
			Message:   "processed request",
			Context: &map[string]interface{}{
				"latency":     end.Sub(start),
				"client_ip":   c.ClientIP(),
				"user_agent":  c.Request.UserAgent(),
				"status_code": c.Writer.Status(),
				"size":        c.Writer.Size(),
				"path":        c.Request.URL.Path,
			},
		}

		loggerInstance.SendLog("user.test", logData)
	}
}
