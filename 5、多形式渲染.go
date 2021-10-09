package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main()  {
	r := gin.Default()

	//json格式输出
	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"html" : "<b>hello,gin框架</b>",
		})
	})

	//原样html输出
	r.GET("/html", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, gin.H{
			"html" : "<b>hello,gin框架</b>",
		})
	})

	//xml输出
	r.GET("/xml", func(c *gin.Context) {
		type Message struct {
			Name string
			Msg string
			Age int
		}
		info := Message{}
		info.Name = "吕振超"
		info.Msg = "hello"
		info.Age = 12
		c.XML(http.StatusOK, info)
	})

	//yml形式
	r.GET("/yaml", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{
			"mes": "gin框架",
			"code": 200,
		})
	})

	r.Run(":9090")
}
