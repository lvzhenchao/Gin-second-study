package middleware

import (
	"Gin-second-study/Gin_App_Api/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err != nil {
				response.Fail(c, nil, fmt.Sprint(err))
				c.Abort()
				return
			}
		}()
	}
}
