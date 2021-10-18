package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pjebs/restgate"
	"net/http"
)

func main() {
	r := gin.Default()

	r.Use(authMiddleware1())
	r.GET("/auth2", func(c *gin.Context) {
		resData := struct {
			Code int	`json:"code"`
			Msg string `json:"msg"`
			Data interface{} `json:"data"`
		}{http.StatusOK, "验证通过", "ok"}
		c.JSON(http.StatusOK, resData)
	})

	r.Run(":9090")
}

var db *sql.DB

func init()  {
	db,_ = SqlDB()
}

func SqlDB()(*sql.DB, error)  {
	DB_TYPE := "mysql"
	DB_HOST := "192.168.33.10"
	DB_PORT := "3306"
	DB_USER := "root"
	DB_NAME := "ginsql"
	DB_PASSWORD := "BspKCZLRZWeHeaTR"
	openString  := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME
	db, err := sql.Open(DB_TYPE, openString)

	return db, err
}


func authMiddleware1() gin.HandlerFunc {//动态配置：用户名、密码配置在数据库中
	return func(c *gin.Context) {
		gate := restgate.New(
			"X-Auth-key",
			"X-Auth-Secret",
			restgate.Database,
			restgate.Config{ //上下key=>secret
				DB: db,
				TableName: "users",
				Key: []string{"keys"},
				Secret: []string{"secrets"},
				HTTPSProtectionOff: true, //允许http访问
			})

		nextCalled := false
		nextAdapter := func(http.ResponseWriter, *http.Request) {
			nextCalled = true
			c.Next()
		}
		gate.ServeHTTP(c.Writer, c.Request, nextAdapter)
		if nextCalled == false {
			c.AbortWithStatus(401)
		}
	}
}
