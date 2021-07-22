package ginm

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/peroperogames/perokit/log"
	"net/http"
	"time"
)

// LoggerMiddle 日志拦截中间件
func LoggerMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		path := c.Request.RequestURI
		clientIP := c.GetHeader("X-Forwarded-For")
		if len(clientIP) == 0 {
			clientIP = c.ClientIP()
		}
		method := c.Request.Method
		statusCode := c.Writer.Status()
		if method != http.MethodHead {
			log.Info(fmt.Sprintf("|STATUS: %d	|Latency: %v	|Client ip: %s	|method: %s	|path: %s	",
				statusCode,
				latency,
				clientIP,
				method,
				path))
		}

	}

}
