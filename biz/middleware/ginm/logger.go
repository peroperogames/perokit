package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/peroperogames/perokit/log"
	"time"
)

// LoggerMiddleWare 日志拦截中间件
func LoggerMiddleWare(c *gin.Context) {
	start := time.Now()
	c.Next()
	end := time.Now()
	latency := end.Sub(start)
	path := c.Request.URL.Path
	clientIP := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	log.Info(fmt.Sprintf("|STATUS: %d	|Latency: %v	|Client ip: %s	|method: %s	|path: %s	",
		statusCode,
		latency,
		clientIP,
		method,
		path))
}
