package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

//示例：https://github.com/gin-contrib/sessions

var sessionName string
var sessionValue string

type MyOption struct {
	sessions.Options
}

func main1() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/session", func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("hello") != "world" {
			session.Set("hello", "world")
			session.Save()
		}

		c.JSON(200, gin.H{"hello": session.Get("hello")})
	})
	r.Run(":8080")
}

func main()  {
	r:=gin.Default()

	//添加中间件, 服务器私有的信息
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/session", func(c *gin.Context) {
		name := c.Query("name")

		if len(name) <= 0 {
			c.JSON(http.StatusBadRequest, "数据错误")
			return
		}

		sessionName = "session_"+name
		sessionValue = "session_value_"+name
		session := sessions.Default(c)
		sessionData := session.Get(sessionName)
		if sessionData != sessionValue {
			//保存session
			session.Set(sessionName, sessionValue)
			o := MyOption{}
			o.Path = "/"
			o.MaxAge = 10
			session.Options(o.Options)
			session.Save()
			c.JSON(http.StatusOK, "首次访问，session已经保存")
			return

		}
		c.JSON(http.StatusOK, "访问成功，您的sessionData:"+sessionData.(string))//断言成字符串类型

	})


	r.Run(":9090")
}
