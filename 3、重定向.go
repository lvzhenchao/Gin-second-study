package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main()  {
	r := gin.Default()

	//一般重定向
	r.GET("/redirect1", func(c *gin.Context) {
		url := "http://www.baidu.com"

		//StatusMovedPermanently, 重定向状态码
		c.Redirect(http.StatusMovedPermanently, url)
	})

	//路由重定向：重定向到具体的路由
	r.GET("/redirect2", func(c *gin.Context) {
		c.Request.URL.Path = "/TestRedirect"
		r.HandleContext(c)
	})
	r.GET("/TestRedirect", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg": "TestRedirect重定向响应成功",
		})
	})

	r.Run(":9090")
}
