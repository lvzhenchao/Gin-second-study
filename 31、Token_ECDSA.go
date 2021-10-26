package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

var (
	err3 error
	eccPrivateKey *ecdsa.PrivateKey
	eccPublicKey *ecdsa.PublicKey
)

type EcdsaUser struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Telephone string `json:"telephone"`
	Password string `json:"password"`
}

type EcdsaClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

func init()  {
	eccPrivateKey, eccPublicKey, err3 = getEcdsaKey(2)
	if err3!= nil {
		panic(err3)
		return
	}
}

func main()  {
	r := gin.Default()

	r.POST("/getToken3", func(c *gin.Context) {
		u := EcdsaUser{}
		err := c.Bind(&u)
		if err != nil {
			c.JSON(http.StatusBadRequest, "参数错误")
			return
		}

		//token分发
		token, err := ecdsaReleaseToken(u)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "生成token错误")
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg": "授权成功",
			"data": token,
		})
	})

	r.POST("/checkToken3", ecdsaTokenMiddle(), func(c *gin.Context) {
		c.JSON(http.StatusOK,"证书有效")
	})

	r.Run(":9090")
}

//验证token
func ecdsaTokenMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := "jiangzhou"
		//获取头部
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") {
			c.JSON(http.StatusUnauthorized,gin.H{
				"code": http.StatusUnauthorized,
				"msg": "前缀错误+token无效",
			})
			c.Abort()
			return
		}

		index := strings.Index(tokenString, auth+":") //找到token前缀对应的位置
		tokenString = tokenString[index+len(auth)+1:] //获取真实的token(开始位置为：索引开始的位置+关键字符的长度+1(:的长度为1))
		claims, err := ecdsaParseToke(tokenString)//解析
		if err != nil {//解析错误或者过期等
			c.AbortWithStatusJSON(http.StatusUnauthorized,err)
			c.Abort()
			return
		}
		claimsValue := claims.(jwt.MapClaims)//断言
		if claimsValue["user_id"] == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,"id不存在")
			c.Abort()
			return
		}

		var u EcdsaUser
		c.Bind(&u)
		if u.Id != claimsValue["user_id"] {
			c.JSON(http.StatusUnauthorized,gin.H{
				"code": http.StatusUnauthorized,
				"msg": "用户不存在",
			})
			c.Abort()
			return
		}
		c.Next()

	}
}

//解析
func ecdsaParseToke(tokenString string) (interface{}, error) {
	myToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		_, ok := token.Method.(*jwt.SigningMethodECDSA) //断言
		if !ok {
			return nil, fmt.Errorf("无效的签名方法：%v", token.Method)
		}
		return eccPublicKey, nil
	})
	if claims,ok := myToken.Claims.(jwt.MapClaims);ok && myToken.Valid  {
		return claims, nil
	}

	return nil, err
}

func ecdsaReleaseToken(u EcdsaUser)(interface{}, error) {
	claims := EcdsaClaims{
		UserId: u.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), //过期时间
			Issuer:    "jiangzhou",                               //发布者
		},
	}
	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedString, err := token.SignedString(eccPrivateKey)
	if err != nil {
		return "", err
	}

	return signedString, nil

}

//ecdsa秘钥生成
func getEcdsaKey(keyType int) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	var err error
	var prk *ecdsa.PrivateKey
	var pub *ecdsa.PublicKey

	var curce elliptic.Curve//椭圆曲线级别类型

	switch keyType {
	case 1:
		curce=elliptic.P224()
	case 2:
		curce=elliptic.P256()
	case 3:
		curce=elliptic.P384()
	case 4:
		curce=elliptic.P521()
	default:
		err = errors.New("输入的签名key类型错误！key取值：\n 1:椭圆曲线224 \n 2:椭圆曲线256 \n 3:椭圆曲线384 \n 4:椭圆曲线521 \n")
		return nil, nil, err

	}

	//生成公私钥对
	prk, err = ecdsa.GenerateKey(curce, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	pub = &prk.PublicKey

	return prk,pub,err
}
