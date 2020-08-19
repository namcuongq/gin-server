package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Data{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
	})
}

func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Data{
		Code: http.StatusOK,
		Data: data,
	})
}

func Success(c *gin.Context, mess ...string) {
	var message = "Success"
	if mess != nil && len(mess) > 0 {
		message = mess[0]
	}

	c.JSON(http.StatusOK, Data{
		Code:    http.StatusOK,
		Message: message,
	})
}

func BadRequest(c *gin.Context, mess ...string) {
	var message = "Bad Request"
	if mess != nil && len(mess) > 0 {
		message = mess[0]
	}

	c.JSON(http.StatusBadRequest, Data{
		Code:    http.StatusBadRequest,
		Message: message,
	})
}

func NotFound(c *gin.Context, mess ...string) {
	var message = "Not Found"
	if mess != nil && len(mess) > 0 {
		message = mess[0]
	}

	c.JSON(http.StatusNotFound, Data{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

func InternalServerError(c *gin.Context, mess ...string) {
	var message = "Internal Server Error"
	if mess != nil && len(mess) > 0 {
		message = mess[0]
	}

	c.JSON(http.StatusInternalServerError, Data{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}
