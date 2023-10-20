package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"net/http"
)

//1、申请证书：可以通过KeyManager来生成测试证书 https://keymanager.org/
//2、证书中间件

type HttpRes struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
}

func main() {
	r := gin.Default()

	r.Use(httpsHandler())
	r.GET("/https_test", func(c *gin.Context) {
		fmt.Println(c.Request.Host)

		c.JSON(http.StatusOK, HttpRes{
			Code:   http.StatusOK,
			Result: "测试成功",
		})
	})

	path := "E:/GoPath/src/CA/"                     //证书路径
	r.RunTLS(":9090", path+"ca.crt", path+"ca.key") //开启HTTPS服务
}

func httpsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLForceHost: true, //只允许https请求
			//SSLHost:"",//http到https的重定向
			STSSeconds:           1536000, //时效
			STSIncludeSubdomains: true,
			STSPreload:           true, //sts预加载
			FrameDeny:            true, //页面不允许frame中展示
			ContentTypeNosniff:   true, //禁用浏览器的类型猜测行为，防止基于MIME类型混淆的攻击
			BrowserXssFilter:     true, //启用XSS保护,并在检查到XSS攻击是，停止渲染页面
		})

		err := secureMiddleware.Process(c.Writer, c.Request)
		//如果不安全，终止拦截
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "数据不安全")
			return
		}

		//判断是否是重定向
		status := c.Writer.Status()
		if status > 300 && status < 399 {
			c.Abort()
			return
		}

		c.Next()

	}
}
