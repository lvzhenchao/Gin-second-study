package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.POST("/post", postMsg)
	r.Run(":9090")
}

func postMsg(c *gin.Context) {
	//name := c.Query("name")

	//name := c.PostForm("name")//获取表单数据
	name := c.DefaultPostForm("name", "lzc") //获取表单数据，并有默认值

	form, b := c.GetPostForm("name") //判断是否存在
	fmt.Println(form, b)

	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"code": http.StatusOK,
	})
}
