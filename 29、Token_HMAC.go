package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

//hash消息认证码
//HS256、HS384、HS512

type HmacUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

type MyClaims struct {
	UserId string
	jwt.StandardClaims
}

var jwtkey = []byte("a_secret_key") //证书签名秘钥（该秘钥非常重要，如果client端有该秘钥，就可以签发证书了）

func main() {
	r := gin.Default()

	//token分发
	r.POST("getToken1", func(c *gin.Context) {

		//{
		//	"id": "20211024",
		//	"name": "lzc",
		//	"telephone": "15910371690",
		//	"password": "000000"
		//}
		var u HmacUser
		c.Bind(&u)
		token, err := hmacReleaseToken(u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "token分发成功",
			"data": token,
		})

	})

	//token认证
	r.POST("/checkToken1", hmacAuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, "证书有效")
	})

	r.Run(":9090")
}

func hmacAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := "jiangzhou" //前缀
		//获取头部
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "前缀错误",
			})
			c.Abort()
			return
		}
		index := strings.Index(tokenString, auth+":")    //找到token前缀对应的位置
		tokenString = tokenString[index+len(auth)+1:]    //获取真实的token(开始位置为：索引开始的位置+关键字符的长度+1(:的长度为1))
		token, claims, err := hamcParseToke(tokenString) //解析

		fmt.Println(err, token.Valid, claims)
		if err != nil || !token.Valid { //解析错误或者过期等
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "证书无效",
			})
			c.Abort()
			return
		}
		var u HmacUser
		c.Bind(&u)
		if u.Id != claims.UserId {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "证书无效",
			})
			c.Abort()
			return
		}
		c.Next()

	}
}

//解析token
func hamcParseToke(tokenString string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{} //cannot unmarshal object into Go value of type jwt.Claims  //地址引用才可以
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, e error) {
		return jwtkey, nil
	})
	return token, claims, err
}

//分发Token
func hmacReleaseToken(u HmacUser) (string, error) {
	expiratiionTime := time.Now().Add(7 * 24 * time.Hour) //截止时间：从当前时刻算起，7天
	claims := &MyClaims{
		UserId: u.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiratiionTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),      //发布者
			Subject:   "user token",           //主题
			Issuer:    "jiangzhou",            //发布者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //生成token
	tokenString, err := token.SignedString(jwtkey)             //签名加密
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
