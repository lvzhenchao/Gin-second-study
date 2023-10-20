package route

import (
	"Gin-second-study/Gin_App_Api/controller"
	"Gin-second-study/Gin_App_Api/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRouter(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware(), middleware.RecoverMiddleware())    //使用中间件
	r.POST("/api/auth/register", controller.Register)                     //注册
	r.POST("/api/auth/login", controller.Login)                           //登录
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info) //再传递数据

	return r
}
