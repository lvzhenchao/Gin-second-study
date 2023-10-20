package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/satori/go.uuid"
	"net/http"
	"unicode/utf8"
)

//通过struct 结构体的tag进行简单校验

//复杂校验：专门库 go-playground/validator

type UserInfo struct {
	Id   string `validate:"uuid" json:"id"`        //设置成uuid,通用识别码
	Name string `validate:"checkName" json:"name"` //自定义了一个规则
	Age  uint8  `validate:"min=0,max=130" json:"age"`
}

//设定全局变量校验
var valildate *validator.Validate

func init() { //对校验进行初始化
	valildate = validator.New()
	valildate.RegisterValidation("checkName", checkNameFunc)
}

func checkNameFunc(f validator.FieldLevel) bool {
	count := utf8.RuneCountInString(f.Field().String()) //一个汉字获取到为长度为1
	if count >= 2 && count <= 12 {
		return true
	}
	return false
}

func main() {
	uuid := uuid.Must(uuid.NewV4(), nil)
	fmt.Println("uuid的值：", uuid)

	r := gin.Default()

	var user UserInfo

	r.POST("/validate", func(c *gin.Context) {
		err := c.Bind(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, "请求参数错误")
			return
		}
		//结构体校验
		err = valildate.Struct(user)
		if err != nil {

			for _, e := range err.(validator.ValidationErrors) {
				fmt.Println("错误的字段：", e.Field())
				fmt.Println("错误的值：", e.Value())
				fmt.Println("错误的Tag：", e.Tag())
			}
			c.JSON(http.StatusBadRequest, "校验数据失败")
			return
		}

		c.JSON(http.StatusOK, "校验成功")
	})

	r.Run(":9090")
}
