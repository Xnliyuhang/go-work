package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		if path == "/users/login" || path == "/users/signup" {
			return
		}
		session := sessions.Default(context)
		if session.Get("userID") == nil {
			//中断,不向下执行
			context.AbortWithStatus(http.StatusUnauthorized)
			context.String(http.StatusOK, "未授权")
			return
		}
	}
}
