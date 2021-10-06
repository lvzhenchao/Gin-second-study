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
	c.String(http.StatusOK, "欢迎您：%s", name)
}
