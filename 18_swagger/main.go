package main

import (
	_"Gin-second-study/18_swagger/docs"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type User struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type Response struct {
	code int `json:"code"`
	Msg string `json:"msg"`
	Data string `json:"data"`
}

func main()  {
	r := gin.Default()

	//swagger中间件主要的作用是：方便前端对接口进行调试。不影响接口的实际功能
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/login", login)
	r.POST("/register", register)

	r.Run(":9090")
}

// @Tags 注册接口
// @Summary 注册
// @Description register
// @Accept json
// @Param  username formData string true "用户名"
// @Param  password formData string true "密码"
// @Success 200 {string} json "{"code":200,"data":"{"name":"username","password":"password"}","msg":"Ok"}"
// @Router /register [post]
func register(c *gin.Context) {
	var user User
	//err := c.Bind(&user)
	err := c.BindQuery(&user)
	if err != nil {
		fmt.Println("绑定错误：", err)
		c.JSON(http.StatusBadRequest, "数据错误")
		return
	}

	res := Response{
		code: http.StatusOK,
		Msg : "注册成功",
		Data: "ok",
	}
	c.JSON(http.StatusOK, res)
}

// @Tags 登录接口
// @Summary 登录
// @Description login
// @Accept json
// @Param  username query string true "用户名"
// @Param  password query string false "密码"
// @Success 200 {string} json "{"code":200,"data":"{"name":"username","password":"password"}","msg":"Ok"}"
// @Router /login [get]
func login(c *gin.Context) {
	userName := c.Query("name")
	pwd := c.Query("pwd")
	fmt.Println(userName, pwd)

	res := Response{}
	res.code = http.StatusOK
	res.Msg = "登录成功"
	res.Data = "ok"
	c.JSON(http.StatusOK,res)
}
