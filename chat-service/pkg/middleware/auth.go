package middleware

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "message": "未授权"})
			c.Abort()
			return
		}

		// 解析token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "message": "无效的token格式"})
			c.Abort()
			return
		}

		// 验证token
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			// 这里应该使用与auth-service相同的密钥
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "message": "无效的token"})
			c.Abort()
			return
		}

		// 从token中获取用户信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "message": "无效的token声明"})
			c.Abort()
			return
		}

		// 将用户ID存入上下文
		userID := int64(claims["user_id"].(float64))
		c.Set("user_id", userID)

		c.Next()
	}
} 