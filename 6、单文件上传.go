package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//附件上传

func main()  {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {

		file, err := c.FormFile("fileName")
		if err != nil {
			c.String(http.StatusBadRequest, "文件上传失败")
		}

		//存储路径地址
		dst := "E:/GoPath/src/"
		c.SaveUploadedFile(file, dst+file.Filename)//存储文件

		c.String(http.StatusOK, fmt.Sprintf("%s 上传文件完成", file.Filename))

	})

	r.Run(":9090")
}
