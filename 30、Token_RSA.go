package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//RAS 秘钥生成工具链接：http://www.metools.info/code/c80.html
//RSA实现Token
type RsaUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

type RasClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

var (
	resPrivateKey  []byte
	resPublicKey   []byte
	err2_1, err2_2 error
)

func init() {
	resPrivateKey, err2_1 = ioutil.ReadFile("E:/GoPath/src/RSA_token/private.pem")
	resPublicKey, err2_2 = ioutil.ReadFile("E:/GoPath/src/RSA_token/public.pem")
	if err2_1 != nil || err2_2 != nil {
		panic(fmt.Sprintf("打开秘钥文件错误：%s,%s", err2_1, err2_2))
		return
	}
}

func main() {
	r := gin.Default()

	r.POST("getToken2", func(c *gin.Context) {
		u := RsaUser{}
		err := c.Bind(&u)
		if err != nil {
			c.JSON(http.StatusBadRequest, "参数错误")
			return
		}
		token, err := resaReleaseToken(u)
		if err != nil {
			c.JSON(http.StatusOK, "生成token错误")
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "token分发成功",
			"data": token,
		})

	})

	//中间件方法加括号会有返回值，不加括号没有返回值
	r.POST("checkToken2", rsaTokenMiddle(), func(c *gin.Context) {
		c.JSON(http.StatusOK, "证书有效")
	})

	r.Run(":9090")
}

func rsaTokenMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := "jiangzhou"
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

		index := strings.Index(tokenString, auth+":") //找到token前缀对应的位置
		tokenString = tokenString[index+len(auth)+1:] //获取真实的token(开始位置为：索引开始的位置+关键字符的长度+1(:的长度为1))
		claims, err := rsaParseToke(tokenString)      //解析

		if err != nil { //解析错误或者过期等
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "证书无效1",
			})
			c.Abort()
			return
		}

		claimsValue := claims.(jwt.MapClaims) //断言
		if claimsValue["user_id"] == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "证书无效2",
			})
			c.Abort()
			return
		}

		var u RsaUser
		c.Bind(&u)
		id := claimsValue["user_id"].(string)
		if u.Id != id {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "用户不存在",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func rsaParseToke(tokenString string) (interface{}, error) {
	pem, err := jwt.ParseRSAPublicKeyFromPEM(resPublicKey)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA) //断言
		if !ok {
			return nil, fmt.Errorf("解析方法错误")
		}
		return pem, err
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

//分发Token
func resaReleaseToken(u RsaUser) (interface{}, error) {
	tokenGen, err := rasJwtTokenGen(u.Id)
	return tokenGen, err
}

//生成Token
func rasJwtTokenGen(id string) (interface{}, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(resPrivateKey)
	if err != nil {
		return nil, err
	}
	claims := RasClaims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), //过期时间
			Issuer:    "jiangzhou",                               //发布者
		},
	}

	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey) //签名加密
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
