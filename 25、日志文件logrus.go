package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"time"
)

//设置最大保存时间，设置日志切割时间间隔

var (
	logFilePath = "./"//文件存储路径
	logFileName = "system.log"
)

func main()  {

	r := gin.Default()

	r.Use(logMiddleware())

	r.GET("/logrus2", func(c *gin.Context) {
		resData := struct {
			Code int `json:"code"`
			Msg string `json:"msg"`
			Data interface{} `json:"data"`
		}{http.StatusOK, "响应成功", "OK"}

		c.JSON(http.StatusOK, resData)
	})


	r.Run(":9090")
}

func logMiddleware() gin.HandlerFunc {
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	file, e := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if e != nil {
		fmt.Println(e)
	}
	//实例化
	logger := logrus.New()
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	//设置输出
	logger.Out=file

	//分割后的文件名称
	logfWriter, _ := rotatelogs.New(
					fileName + ".%Y%m%d.log",
					rotatelogs.WithLinkName(fileName),//生成软链接，指向最新日志文件
					rotatelogs.WithMaxAge(7*24 * time.Hour),//设置最大保存时间（7天）
					rotatelogs.WithRotationTime(time.Hour),//日志切割时间（1天）
				)

	//hook机制设置
	WriterMaps := lfshook.WriterMap{
		logrus.InfoLevel:  logfWriter,
		logrus.FatalLevel: logfWriter,
		logrus.DebugLevel:  logfWriter,
		logrus.WarnLevel:  logfWriter,
		logrus.ErrorLevel:  logfWriter,
		logrus.PanicLevel:  logfWriter,
	}

	//给logrus添加hook
	logger.AddHook(lfshook.NewHook(WriterMaps, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))

	return func(c *gin.Context) {
		logger.WithFields(logrus.Fields{//里面的值都可以自己定义
			"status_code": c.Writer.Status(),
			"url": c.Request.RequestURI,
			"method": c.Request.Method,
			"params": c.Query("name"),
			"IP":c.ClientIP(),
		}).Info()//显示日志
	}

}
