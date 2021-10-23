package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

//特别注意：结构体名称为：Product, 创建的表的名称为：Products

type Product struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Number         string    `grom:"unique" json:"number"`          //类别
	Category       string    `gorm:"type:varchar(256);not null" json:"category"`        //分类
	Name           string    `gorm:"type:varchar(20);not null" json:"name"`            //商品命称
	MadeIn         string    `gorm:"type:varchar(128);not null" json:"made_in"`         //生产地
	ProductionTime time.Time `json:"production_time"` //生产时间
}

//应答客户端的请求
type GormResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

var gormDB *gorm.DB
var gormResponse GormResponse

func init()  {
	//1、打开数据库
	sqlStr := "root:BspKCZLRZWeHeaTR@tcp(192.168.33.10:3306)/ginsql?charset=utf8mb4&parseTime=true&&loc=Local"
	var err error
	gormDB, err = gorm.Open(mysql.Open(sqlStr), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接错误")
		return
	}

}

func main() {
	r := gin.Default()

	r.POST("gorm/insert", gormInsertData)
	r.GET("gorm/get", gormGetData)
	r.GET("gorm/mulget", gormGetMulData)
	r.PUT("gorm/update", gormUpdate)//修改数据
	r.DELETE("gorm/delete", gormDelete)//删除数据

	r.Run(":9090")
}

func gormDelete(c *gin.Context) {
	number := c.Query("number")
	//1、先查询
	var count int64
	gormDB.Model(&Product{}).Where("number=?", number).Count(&count)
	if count <= 0 {
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "查询错误"
		gormResponse.Data = "error"
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	//2、删除
	tx := gormDB.Where("number=?", number).Delete(&Product{})
	if tx.RowsAffected >0 {
		gormResponse.Code = http.StatusOK
		gormResponse.Message = "删除成功"
		gormResponse.Data = "ok"
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	fmt.Printf("删除错误：err:%v\n", tx)
	gormResponse.Code = http.StatusBadRequest
	gormResponse.Message = "删除错误"
	gormResponse.Data = tx
	c.JSON(http.StatusOK, gormResponse)
}

func gormUpdate(c *gin.Context) {
	var p Product
	err := c.Bind(&p)
	if err != nil{
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "参数错误"
		gormResponse.Data = err
		c.JSON(http.StatusOK, gormResponse)
		return
	}

	//1、先查询
	var count int64
	gormDB.Model(&Product{}).Where("number=?", p.Number).Count(&count)
	if count <= 0 {
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "查询错误"
		gormResponse.Data = "error"
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	//2、更新
	tx := gormDB.Model(&Product{}).Where("number=?", p.Number).Updates(&p)
	if tx.RowsAffected >0 {
		gormResponse.Code = http.StatusOK
		gormResponse.Message = "更新成功"
		gormResponse.Data = "ok"
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	fmt.Printf("更新错误：err:%v\n", tx)
	gormResponse.Code = http.StatusBadRequest
	gormResponse.Message = "更新错误"
	gormResponse.Data = tx
	c.JSON(http.StatusOK, gormResponse)
}

func gormGetMulData(c *gin.Context) {
	category := c.Query("category")
	products := make([]Product, 10)
	tx := gormDB.Where("category=?", category).Find(&products).Limit(10)
	if tx.Error != nil {
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "查询错误"
		gormResponse.Data = tx.Error
		c.JSON(http.StatusOK, gormResponse)
		return
	}

	gormResponse.Code = http.StatusOK
	gormResponse.Message = "查询成功"
	gormResponse.Data = products
	c.JSON(http.StatusOK, gormResponse)
}

func gormGetData(c *gin.Context) {
	number := c.Query("number")
	product := Product{}
	tx := gormDB.Where("number=?", number).First(&product)
	if tx.Error != nil {
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "查询错误"
		gormResponse.Data = tx.Error
		c.JSON(http.StatusOK, gormResponse)
		return
	}

	gormResponse.Code = http.StatusOK
	gormResponse.Message = "查询成功"
	gormResponse.Data = product
	c.JSON(http.StatusOK, gormResponse)
}

func gormInsertData(c *gin.Context) {
	var p Product
	err := c.Bind(&p)
	if err != nil{
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "参数错误"
		gormResponse.Data = "error"
		c.JSON(http.StatusOK, gormResponse)
		return
	}

	tx := gormDB.Create(&p)
	if tx.RowsAffected > 0 {
		gormResponse.Code = http.StatusOK
		gormResponse.Message = "写入成功"
		gormResponse.Data = "OK"
		c.JSON(http.StatusOK, gormResponse)
		return
	}

	fmt.Printf("插入错误：err:%v\n", tx)
	gormResponse.Code = http.StatusBadRequest
	gormResponse.Message = "写入失败"
	gormResponse.Data = tx
	c.JSON(http.StatusOK, gormResponse)
}
