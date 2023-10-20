package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

//多服务器程序运行

var g errgroup.Group

func main() {
	//服务器1：
	server01 := http.Server{
		Addr:              ":9091",
		Handler:           roter01(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	//服务器2：
	server02 := http.Server{
		Addr:              ":9092",
		Handler:           roter02(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	//开启对应的服务
	g.Go(func() error {
		return server01.ListenAndServe()
	})
	g.Go(func() error {
		return server02.ListenAndServe()
	})

	//阻塞一下
	if err := g.Wait(); err != nil {
		fmt.Println("执行失败", err)
	}
}

func roter01() http.Handler {
	r1 := gin.Default()

	r1.GET("/MyServer01", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "服务器程序1的响应",
		},
		)
	})
	return r1
}

func roter02() http.Handler {
	r2 := gin.Default()

	r2.GET("/MyServer02", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "服务器程序2的响应",
		},
		)
	})
	return r2
}
