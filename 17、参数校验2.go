package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

//嵌套结构体的校验：dive关键字，代表进入到嵌套结构体进行判断

type ValUser struct {
	Name    string       `validate:"required" json:"name"`
	Age     uint8        `validate:"gte=0,lte=130" json:"age"`    // 0<=Age<=130
	Email   string       `validate:"required,email" json:"email"` //非空，email格式
	Address []ValAddress `validate:"dive" json:"address"`         //嵌套结构体验证
}

type ValAddress struct {
	Province string `validate:"required" json:"province"`
	City     string `validate:"required" json:"city"`
	Phone    string `validate:"numeric,len=11" json:"phone"`
}

var validate *validator.Validate

func init() {
	validate = validator.New() //初始化（复制）
}

func testData(c *gin.Context) {
	address := ValAddress{
		Province: "北京",
		City:     "大兴",
		Phone:    "15910371699",
	}
	user := ValUser{
		Name:    "lzc",
		Age:     12,
		Email:   "123@qq.com",
		Address: []ValAddress{address, address},
	}
	c.JSON(http.StatusOK, user)
}

func main() {
	r := gin.Default()

	var user ValUser
	r.POST("/validate", func(c *gin.Context) {

		//testData(c)//生成用于postman的测试数据

		err := c.Bind(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, "参数错误")
			return
		}

		//执行参数的校验
		res := validateUser(user)
		if res {
			c.JSON(http.StatusOK, "校验成功")
		} else {
			c.JSON(http.StatusBadRequest, "校验失败")
		}

	})

	r.Run(":9090")
}

func validateUser(u ValUser) bool {
	err := validate.Struct(u)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println("错误的字段：", e.Field())
			fmt.Println("错误的值：", e.Value())
			fmt.Println("错误的Tag：", e.Tag())
		}
		return false

	}

	return true
}
