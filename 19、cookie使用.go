package main

import (
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//主要是使用中间件使用cookie
var cookieName string
var cookieValue string


func main()  {
	r := gin.Default()

	//cookie中间件,保存cookie
	r.Use(CookieAuth())

	r.GET("/cookie", func(c *gin.Context) {

		name := c.Query("name")

		if len(name) <= 0 {
			c.JSON(http.StatusBadRequest, "数据错误")
			return
		}

		cookieName = "cookie_"+name
		cookieValue = hex.EncodeToString([]byte(cookieName+"value"))

		val, _ := c.Cookie(cookieName)
		if val == "" {
			c.String(http.StatusOK, "%scookie已下发，下次登录有效", cookieName)
			return
		}
		c.String(http.StatusOK, "验证成功，cookie值为%s", val)


	})

	r.Run(":9090")
}

func CookieAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Cookie(cookieName)
		if val == "" {
			c.SetCookie(cookieName, cookieValue, 3600, "/", "localhost", true, true)
			fmt.Println("cookie已经保存成功")
		}
	}
}
