package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main()  {
	r := gin.Default()

	r.GET("/get", getMsg)

	r.Run(":9090")
}

func getMsg(c *gin.Context)  {
	name := c.Query("name")
	//返回字符串类型
	//c.String(http.StatusOK, "欢迎您：%s", name)

	//返回json数据类型
	c.JSON(http.StatusOK, gin.H{
		"code" : http.StatusOK,
		"msg" : "返回信息",
		"data": name,
	})
}
