package common

import (
	"Gin-second-study/Gin_App_Api/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("a_secret_key") //证书签名秘钥（该秘钥非常重要，如果client端有该秘钥，就可以签发证书了）

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

//分发证书
func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) //截止时间：从当前时刻算起，7天
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),     //发布时间
			Issuer:    "jiangzhou",           //发布者
			Subject:   "user token",          //主题
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //生成token
	tokenString, err := token.SignedString(jwtKey)             //签名

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//解析证书
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}
