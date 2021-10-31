package controller

import (
	"Gin-second-study/Gin_App_Api/common"
	"Gin-second-study/Gin_App_Api/model"
	"Gin-second-study/Gin_App_Api/response"
	"Gin-second-study/Gin_App_Api/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

func Register(ctx *gin.Context)  {
	var requestUser model.User
	ctx.Bind(&requestUser)

	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity,422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity,422, nil, "密码不能少于6位")
		return
	}
	if name == "" {
		name = util.RandomString(10)
	}
	if isTelephoneExist(common.DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}

	//创建用户
	//返回密码的hash值（对用户密码进行二次处理，防止系统管理人员利用）
	hassPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}

	newUser := model.User{
		Name: name,
		Telephone:telephone,
		Password:string(hassPassword),//byte转为字符串类型
	}
	common.DB.Create(&newUser)//新增用户

	//分发token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg" : "系统异常",
		})
		return
	}

	//返回前端
	response.Success(ctx, gin.H{
		"token": token,
	}, "注册成功")

}

func Login(ctx *gin.Context)  {
	var requestUser model.User
	ctx.Bind(&requestUser)

	//name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity,422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity,422, nil, "密码不能少于6位")
		return
	}

	//依据手机号，查询用户注册的数据记录
	var user model.User
	common.DB.Where("telephone=?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	//判断密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//分发token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg" : "系统异常",
		})
		return
	}

	//返回前端
	response.Success(ctx, gin.H{
		"token": token,
	}, "登录成功")
}

func Info(ctx *gin.Context)  {
	user, _ := ctx.Get("user")//空接口
	response.Success(ctx, gin.H{
		"user": response.ToUserDto(user.(model.User)),//将空接口断言成结构体
	}, "响应成功")
	return
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone=?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
