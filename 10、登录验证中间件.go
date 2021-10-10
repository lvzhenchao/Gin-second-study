package main

//浏览器访问：login；一次登录就可以；除非关掉浏览器

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main()  {
	r := gin.Default()

	r.Use(AuthMiddleware())

	r.GET("/login", func(c *gin.Context) {

		//这里获取登录的用户，它是由basicAuth中间件获取的
		user := c.MustGet(gin.AuthUserKey).(string)//断言

		c.JSON(http.StatusOK, "登录成功! "+"欢迎:"+user)
	})

	r.Run(":9090")
}

func AuthMiddleware() gin.HandlerFunc {
	//初始化用户
	accounts := gin.Accounts{//静态的账号
		"admin" : "adminpw",
		"system" : "systempw",
	}

	//动态添加账号
	accounts["go"] = "123456"
	accounts["gin"] = "gin123"

	//将用户添加到登录中间件
	auth := gin.BasicAuth(accounts)

	return auth
}