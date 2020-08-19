package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjk/dailyrotate"
)

func Logger(logFile *dailyrotate.File) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		logAccess := fmt.Sprintf("%s|%v|%v|%v|%v|%v|%v\n",
			start.Format("15:04:05 02-01-2006"),
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.Method,
			path,
			c.Request.UserAgent())
		logFile.Write([]byte(logAccess))
	}
}
