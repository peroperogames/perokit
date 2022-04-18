package ginm

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/peroperogames/perokit/log"
	"net/http"
	"time"
)

// LoggerMiddle 日志拦截中间件
func LoggerMiddle(logger log.Logger) gin.HandlerFunc {
	helper := log.NewHelper(log.With(logger, "middle", "middle/LoggerMiddle"))
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		path := c.Request.RequestURI
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		clientIP := c.GetHeader("X-Forwarded-For")
		if len(clientIP) == 0 {
			clientIP = c.ClientIP()
		}
		method := c.Request.Method
		statusCode := c.Writer.Status()
		if method != http.MethodHead && clientIP != "127.0.0.1" {
			helper.Info(fmt.Sprintf("|STATUS: %3d|Latency: %13v|Client ip: %15s|method: %s|path: %s",
				statusCode,
				latency,
				clientIP,
				method,
				path))
		}

	}

}
