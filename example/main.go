package main

import (
	"fmt"
	"gin-server/jwt"
	"gin-server/middleware"
	"gin-server/response"
	"gin-server/server"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := server.New(server.ENV_DEVELOPMENT)
	router.LoadStatic("/static/", "<path to static folder>")
	router.LoadTemplate("<path to template folder>")
	router.CORS(middleware.CORSConfig{
		AllowOrigin:  "<domain>",
		MaxAge:       "8600",
		AllowMethods: "GET, POST",
		AllowHeaders: "Content-Type",
	})
	router.AllowHosts("abc.com", "abc.vn")
	router.SecureHeader()
	router.FuckBot()
	router.AccessLogDaily("<folder-access-log>")

	authen, err := jwt.New(&jwt.JWT{
		SecretKey:     "secret-key",
		ExpiredHour:   1,        //deault 1 hour
		TokenHeadName: "Bearer", // TokenHeadName is a string in the header. Default value is ""
		Authenticator: func(c *gin.Context) (map[string]interface{}, error) {
			var loginVals map[string]string
			if err := c.ShouldBind(&loginVals); err != nil {
				return nil, fmt.Errorf("error username or password missing")
			}

			if loginVals["username"] != "admin" || loginVals["password"] != "admin" {
				return nil, fmt.Errorf("authentication failed")
			}

			var data = map[string]interface{}{
				"id": 1,
			}

			return data, nil
		},
		Verification: func(c *gin.Context, data map[string]interface{}) (bool, error) {
			if fmt.Sprintf("%v", data["username"]) == "admin" {
				return true, nil
			}

			return false, nil
		},
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	api := router.Group("/api")
	api.POST("/login", authen.LoginHandler)
	api.GET("/refresh", authen.RefreshToken)

	api.Use(authen.TokenAuthMiddleware())
	api.GET("/home", func(c *gin.Context) {
		response.SuccessWithData(c, map[string]interface{}{
			"hello": "admin",
		})
	})

	router.Run()
}
