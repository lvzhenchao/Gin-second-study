package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//同步：调用开始，必须等到方法调用返回后，才能继续后续的行为
//异步：调用更像一个消息传递，异步方法通常会在另外一个go程（协程）中进行，不会阻碍调用者的工作

//特别注意：可以在中间件或处理程序中启动新的Go协程
//特别注意：需要使用上下文的副本【子进程的上下文是对主进程上下文的拷贝，不会影响主进程的上下文】

func main()  {
	r := gin.Default()

	r.GET("/sync", func(c *gin.Context) {
		sync(c)
		c.JSON(http.StatusOK, ">>>主进程同步已经执行<<<")
	})

	r.GET("/async", func(c *gin.Context) {
		for i:= 0; i<6; i++ {

			cCp := c.Copy()//特别注意：需要使用上下文的副本
			go async(cCp, i)

		}
		c.JSON(http.StatusOK, ">>>主进程异步已经执行<<<")
	})

	r.Run(":9090")
}

func async(cp *gin.Context, i int) {
	fmt.Println("第"+strconv.Itoa(i)+"个go进程开始执行："+cp.Request.URL.Path)
	time.Sleep(time.Second*3)
	fmt.Println("第"+strconv.Itoa(i)+"个go进程结束执行："+cp.Request.URL.Path)
}



func sync(c *gin.Context) {
	println("开始执行同步任务："+c.Request.URL.Path)
	time.Sleep(time.Second*3)
	println("同步任务执行完毕！")
}
