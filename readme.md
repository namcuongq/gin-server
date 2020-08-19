# Gin Server
###
Công cụ để tạo WEB Server sử dụng [Gin Framework](https://github.com/gin-gonic/gin) cùng mới nhiều tính năng mở rộng. Dùng để xây dựng WEB Server một cách nhanh chóng và bảo mật.
###
## Danh sách tính năng

+ **LoadStatic** - Load Static File sử dụng [Static](github.com/gin-contrib/static)
+ **LoadTemplate** - Load html template
+ **CORS** - Cài đặt CORS
+ **AllowHosts** - Chỉ cho phép truy cập từ danh sách các hosts
+ **SecureHeader** - Chèn các header giúp bảo mật website 
+ **FuckBot** - Chặn Bot truy cập website
+ **AccessLogDaily** - Ghi log access theo từng ngày 
+ **JWT** - Xác thực bằng JWT 
+ **Limit** - Giới hạn số lượt truy cập của 1 IP (Đang cập nhât)

### Ví dụ

```go
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

```
