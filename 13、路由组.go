package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//路由组可以方便的对路由组进行分组和有效的分类，使路由对应的代码易于阅读

type ResGroup struct {
	Data string
	Path string
}

func main()  {
	router := gin.Default()
	{

		v1 := router.Group("/v1")//路由分组：1级路径
		r := v1.Group("/user")//路由分组：2级路径
		r.GET("/login", login)//请求路径：/v1/user/login

		r2 := r.Group("/showInfo")
		r2.GET("/abstract", abstract)//请求路径：/v1/user/showInfo/abstract
		r2.GET("/detail", detail)//请求路径：/v1/user/showInfo/detail
	}

	v2 := router.Group("/v2")
	{
		v2.GET("/other", other)//v2/other
	}

	router.Run(":9090")
}

func other(c *gin.Context) {
	c.JSON(http.StatusOK, ResGroup{"other", c.Request.URL.Path})
}

func detail(c *gin.Context) {
	c.JSON(http.StatusOK, ResGroup{"detail", c.Request.URL.Path})
}

func abstract(c *gin.Context) {
	c.JSON(http.StatusOK, ResGroup{"abstract", c.Request.URL.Path})
}

func login(c *gin.Context) {
	c.JSON(http.StatusOK, ResGroup{"login", c.Request.URL.Path})
}
