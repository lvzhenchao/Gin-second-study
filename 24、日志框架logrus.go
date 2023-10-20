package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

//logrus是一个结构化、可插拔的Go日志框架；自定义插件功能，有text和json输出格式
//支持field机制和可扩展的hook机制，允许用户通过hook方式将日志分发到任意地方
//docker、prometheus等都是使用logrus记录

var log = logrus.New() //创建一个log示例

func initLogrus() error {
	log.Formatter = &logrus.JSONFormatter{}                                            //设置json格式的日志
	file, e := os.OpenFile("./gin_log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) //创建一个log日志文件
	if e != nil {
		fmt.Println("创建文件|打开文件失败")
		return e
	}

	log.Out = file               //设置log的默认文件输出
	gin.SetMode(gin.ReleaseMode) //设定发布的版本
	gin.DefaultWriter = log.Out  //将gin框架的默认日志信息也输出到里面
	log.Level = logrus.InfoLevel //设置日志级别
	return nil
}

func main() {

	err := initLogrus()
	if err != nil {
		fmt.Println(err)
		return
	}

	r := gin.Default()

	r.GET("/logrus", func(c *gin.Context) {
		log.WithFields(logrus.Fields{ //里面的值都可以自己定义
			"url":    c.Request.RequestURI,
			"method": c.Request.Method,
			"params": c.Query("name"),
			"IP":     c.ClientIP(),
		}).Info() //显示日志

		resData := struct {
			Code int         `json:"code"`
			Msg  string      `json:"msg"`
			Data interface{} `json:"data"`
		}{http.StatusOK, "响应成功", "OK"}

		c.JSON(http.StatusOK, resData)
	})

	r.Run(":9090")
}
