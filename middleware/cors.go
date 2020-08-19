package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	AllowOrigin  string
	MaxAge       string
	AllowMethods string
	AllowHeaders string
}

func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
	if len(config.AllowOrigin) < 1 {
		config.AllowOrigin = "*"
	}

	if len(config.MaxAge) < 1 {
		config.MaxAge = "86400"
	}

	if len(config.AllowHeaders) < 1 {
		config.AllowHeaders = "POST, GET, PUT, DELETE"
	}

	if len(config.AllowMethods) < 1 {
		config.AllowMethods = "Authorization, Content-Type"
	}

	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Writer.Header().Set("Access-Control-Allow-Origin", config.AllowOrigin)
			c.Writer.Header().Set("Access-Control-Max-Age", config.MaxAge)
			c.Writer.Header().Set("Access-Control-Allow-Methods", config.AllowHeaders)
			c.Writer.Header().Set("Access-Control-Allow-Headers", config.AllowMethods)
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
			// if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
