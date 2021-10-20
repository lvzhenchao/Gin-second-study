package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

var sqlDb *sql.DB //数据库连接db
var sqlResponse SqlResponse //响应client的数据

func init()  {
	//1、打开数据库
	//parseTime 时间格式转换
	sqlStr := "root:BspKCZLRZWeHeaTR@tcp(192.168.33.10:3306)/ginsql?charset=utf8&parseTime=true&&loc=Local"
	var err error
	sqlDb, err = sql.Open("mysql", sqlStr)
	if err != nil {
		fmt.Println("数据库打开出现问题：", err)
		return
	}

	//2、测试与数据库建立连接（校验连接是否正确）
	err = sqlDb.Ping()
	if err != nil {
		fmt.Println("数据库连接出现了问题：", err)
		return
	}
}

//客户端提交的数据
type SqlUser struct {
	Name string `json:"name"`
	Age int	`json:"age"`
	Address string `json:"address"`
}

//服务端的响应
type SqlResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func main()  {
	r := gin.Default()

	r.POST("sql/insert", insertData)//写数据
	r.GET("sql/get", getData)//查询单条数据

	r.Run(":9090")
}

func getData(c *gin.Context) {
	name := c.Query("name")
	sqlStr := "select age, address from user where name=? "
	var u SqlUser
	//查询单条，并扫描到结构体
	err := sqlDb.QueryRow(sqlStr, name).Scan(&u.Age, &u.Address)
	if err != nil {
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "查询错误"
		sqlResponse.Data = "error"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}

	u.Name = name
	sqlResponse.Code = http.StatusOK
	sqlResponse.Message = "读取成功"
	sqlResponse.Data = u
	c.JSON(http.StatusOK, sqlResponse)
}

func insertData(c *gin.Context) {
	var u SqlUser
	err := c.Bind(&u)
	if err != nil {
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "参数错误"
		sqlResponse.Data = "error"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}

	fmt.Println(u)
	sqlStr := "insert into user (name, age, address) values (?,?,?)"
	result, err := sqlDb.Exec(sqlStr, u.Name, u.Age, u.Address)
	if err != nil {
		fmt.Printf("插入错误：err:%v\n", err)
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "写入失败"
		sqlResponse.Data = "error"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}


	sqlResponse.Code = http.StatusOK
	sqlResponse.Message = "写入成功"
	sqlResponse.Data = "OK"
	c.JSON(http.StatusOK, sqlResponse)
	fmt.Println(result.LastInsertId())//打印结果




}























