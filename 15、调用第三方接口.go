package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

//将14课 go build成exe运行；然后去请求调用
//记得关掉360

//调用第三方接口的请求数据
type UserAPI struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

//第三方返回的结果
type TempData struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

//客户端提交的数据
type ClientRequest struct {
	UserName string      `json:"user_name"`
	Password string      `json:"password"`
	other    interface{} `json:"other"`
}

//返回客户端的数据
type ClientResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func main() {
	//单独请求测试
	//testAPI()

	r := gin.Default()

	r.POST("/getOtherAPI", getOtherAPI)

	r.Run(":9091")
}

func getOtherAPI(c *gin.Context) {
	var requestData ClientRequest
	var response ClientResponse
	err := c.Bind(&requestData)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Msg = "请求参数错误"
		response.Data = err
		c.JSON(http.StatusBadRequest, response)
		return
	}

	url := "http://127.0.0.1:9090/login"

	user := UserAPI{requestData.UserName, requestData.Password}

	data, err := getRestfulAPI(url, user, "application/json")
	fmt.Println(data, err)

	//json反序列化
	var temp TempData
	json.Unmarshal(data, &temp)

	fmt.Println(temp.Msg, temp.Data)

	response.Code = http.StatusOK
	response.Msg = "请求成功"
	response.Data = temp
	c.JSON(http.StatusBadRequest, response)

}

func testAPI() {
	url := "http://127.0.0.1:9090/login"
	user := UserAPI{"user", "1232"}

	data, err := getRestfulAPI(url, user, "application/json")
	fmt.Println(data, err)

	//json反序列化
	var temp TempData
	json.Unmarshal(data, &temp)

	fmt.Println(temp.Msg, temp.Data)
}

func getRestfulAPI(url string, data interface{}, contentType string) ([]byte, error) {
	//创建调用api接口的client
	client := http.Client{Timeout: 5 * time.Second}

	//转换成json
	jsonStr, _ := json.Marshal(data)

	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("调用API接口出现错误")
		return nil, err
	}

	res, err := ioutil.ReadAll(resp.Body)
	return res, err

}
