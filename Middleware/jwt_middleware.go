package Middleware

import (
	"ddlBackend/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// JWTAuthMiddleware JWT验证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// Authorization 的header由"Bearer " + token组成
		bearerLength := len("Bearer ")

		authHeader := context.GetHeader("Authorization")

		if len(authHeader) < bearerLength {
			// 如果令牌的长度不够
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "The authorization header is incorrect",
			})
			context.Abort()
			return
		}

		token := strings.TrimSpace(authHeader[bearerLength:])
		claims, err := tool.ParseJWTToken(token)
		if err != nil {
			// 解析令牌中遇到错误
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			context.Abort()
		} else if time.Now().Unix() > claims.ExpiresAt.Unix() {
			// 令牌过期
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "The token has been expired",
			})
			context.Abort()
		} else {
			// 没有问题就把信息记录在context中
			context.Set("Claims", *claims)
		}
		return
	}
}
