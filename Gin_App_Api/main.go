package main

import (
	"Gin-second-study/Gin_App_Api/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)

//直接命令行 go run main.go   输出E:\GoPath\src\Gin-second-study\Gin_App_Api
//单击main左边的 run go build main.go  输出E:\GoPath\src\Gin-second-study
//注意区别

func main()  {
	IninConfig()//加载配置项
	common.InitDB()//初始化数据库（只初始化一次）
	r :=gin.Default()

	port := viper.GetString("server.port")
	r.Run(":"+port)
}

func IninConfig() {
	workDir, _ := os.Getwd()//获取目录对应的路径
	viper.SetConfigName("application")//配置文件名
	viper.SetConfigType("yml")//配置文件类型
	viper.AddConfigPath(workDir+"/config")//执行go run对应的路径配置
	fmt.Println(workDir)
	//加载
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}
