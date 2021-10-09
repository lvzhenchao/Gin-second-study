package main

import "github.com/gin-gonic/gin"

//客户端请求内容为 视频、音频、图片等文件

func main()  {
	r := gin.Default()

	r.GET("/file", fileServer)

	r.Run(":9090")
}

func fileServer(c *gin.Context)  {
	path := "E:/GoPath/src/"
	fileName := path + c.Query("name")
	c.File(fileName)
}
