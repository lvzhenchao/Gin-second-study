package route

import (
	"Gin-second-study/Gin_App_Api/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRouter(r *gin.Engine)*gin.Engine  {
	r.Use(middleware.CORSMiddleware(), middleware.RecoverMiddleware())//使用中间件
	r.POST("/api/auth/register",)//注册
	r.POST("/api/auth/login",)//登录
	r.GET("/api/auth/info", middleware.AuthMiddleware())//再传递数据
}


