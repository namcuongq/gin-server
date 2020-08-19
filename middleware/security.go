package middleware

import (
	"gin-server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SecureHeader() gin.HandlerFunc {
	return func(cc *gin.Context) {
		cc.Header("X-Powered-By", "PHP/7.2.24")
		cc.Header("Server", "nginx/1.17.0")
		cc.Header("X-Content-Type-Options", "nosniff")
		cc.Header("X-Frame-Options", "DENY")
		cc.Header("X-XSS-Protection", "1; mode=block")
		cc.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		cc.Header("Cache-Control", "private, max-age=0")
		cc.Header("Pragma", "no-cache")
		cc.Next()
	}
}

func CheckHost(domains ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.IsStringInArrray(c.Request.Host, domains) {
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}
		c.Next()
	}
}

func LimitRequest(num int) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func BlockUserAgentMalicious() gin.HandlerFunc {
	return func(c *gin.Context) {
		ua := utils.ParseUserAgent(c.Request.UserAgent())
		if ua.OS == "" {
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}
		c.Next()
	}
}
