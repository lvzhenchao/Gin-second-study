package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//将客户端提交的json数据与server对应的对象（实体或结构体）进行关联
//bind(类似序列化和反序列化)
//请求参数json对应的key就是结构体对应的字段

type Login struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remark string `json:"remark"`
}

func main()  {
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		var login Login
		err := c.Bind(&login)

		fmt.Println("绑定的数据：", login)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "绑定失败",
				"data": err.Error(),
			})
			return
		}
		if login.UserName == "user" && login.Password == "123" {
			c.JSON(http.StatusOK, gin.H{
				"msg": "登录成功",
				"data": "ok",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "登录失败",
			"data": "error",
		})
		return

	})

	r.Run(":9090")
}
