package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"net/http"
	"time"
)

var x *xorm.Engine
var xormResponse XormResponse

//定义结构体（xorm支持双向映射）：没有表，会进行创建
type Stu struct {
	Id int64 `xorm:"pk autoincr" json:"id"`
	StuNum string `xorm:"unique" json:"stu_num"`
	Name string `json:"name"`
	Age int `json:"age"`
	Created time.Time `xorm:"created" json:"created"`
	Updated time.Time `xorm:"updated" json:"updated"`
}

//应答客户端的请求
type XormResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}


//注意变量设定 runtime error: invalid memory address or nil pointer dereference

func init()  {
	//1、打开数据库
	sqlStr := "root:BspKCZLRZWeHeaTR@tcp(192.168.33.10:3306)/ginsql?charset=utf8&parseTime=true&&loc=Local"
	var err error
	x, err = xorm.NewEngine("mysql", sqlStr)
	if err != nil {
		fmt.Println("数据库连接错误")
	}

	//2、创建或者同步表（名称stu）
	err = x.Sync(new(Stu))
	if err != nil {
		fmt.Println("数据表同步失败：", err)
	}
}


func main()  {
	r:=gin.Default()

	//数据库的增删改查
	r.POST("xorm/insert", xormInsertData)//新增数据
	r.GET("xorm/get", xormGetData)//获取单条数据


	r.Run(":9090")
}

func xormGetData(c *gin.Context) {
	stuNum := c.Query("stu_num")
	var stus []Stu
	err := x.Where("stu_num=?", stuNum).Find(&stus)
	if err != nil {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "查询错误"
		xormResponse.Data = "error"
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	xormResponse.Code = http.StatusOK
	xormResponse.Message = "查询成功"
	xormResponse.Data = stus
	c.JSON(http.StatusOK, xormResponse)

}

func xormInsertData(c *gin.Context) {
	var s Stu
	err := c.Bind(&s)
	if err != nil{
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "参数错误"
		xormResponse.Data = "error"
		c.JSON(http.StatusOK, xormResponse)
		return
	}

	fmt.Println(s)
	affected, err := x.Insert(s)
	if err!= nil || affected <= 0 {
		fmt.Printf("插入错误：err:%v\n", err)
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "写入失败"
		xormResponse.Data = err
		c.JSON(http.StatusOK, xormResponse)
		return
	}

	xormResponse.Code = http.StatusOK
	xormResponse.Message = "写入成功"
	xormResponse.Data = "OK"
	c.JSON(http.StatusOK, xormResponse)
	fmt.Println(affected)//打印结果

}
