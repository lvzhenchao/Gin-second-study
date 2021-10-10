package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

//对clinet请求的路由数据进行预处理（数据加载、请求验证等）

func main()  {
	r := gin.Default()

	r.Use(Middleware())

	r.GET("/middleware", func(c *gin.Context) {
		fmt.Println("服务端开始执行...")
		name := c.Query("name")
		ageStr := c.Query("age")
		age, _ := strconv.Atoi(ageStr)//字符串转整形
		log.Println(name, age)
		res := struct {//加标签，大写转小写
			Name string `json:"name"`
			Age int	`json:"age"`
		}{name, age}

		c.JSON(http.StatusOK, res)
	})

	r.Run(":9090")
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("中间件开始执行...")

		name := c.Query("name")
		ageStr := c.Query("age")

		age, err := strconv.Atoi(ageStr)//字符串转整形
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "数据错误，年龄不是整数")
		}
		if age < 0 || age > 100 {
			c.AbortWithStatusJSON(http.StatusBadRequest, "数据错误， 年龄数据错误")
		}

		if len(name) <= 6 || len(name) >= 12 {
			c.AbortWithStatusJSON(http.StatusBadRequest, "数据错误， 用户名只能是6~12位")
		}

		c.Next()//执行后续操作

		fmt.Println(name, age)
	}
}

