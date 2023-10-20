package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		//获取对应的form
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, "文件上传失败")
		}

		files := form.File["file_key"] //所有的文件名称

		//存储路径地址
		dst := "E:/GoPath/src/"
		for _, file := range files {
			err = c.SaveUploadedFile(file, dst+file.Filename) //存储文件
			if err != nil {
				c.String(http.StatusBadRequest, "文件上传失败")
			}
		}

		//返回信息
		c.String(http.StatusOK, fmt.Sprintf("%d 上传文件完成", len(files)))
	})

	r.Run(":9090")
}
