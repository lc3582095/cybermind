package middleware

import (
    "net/http"
    "strings"
    "cybermind/auth-service/pkg/jwt"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "message": "authorization header is required"})
            c.Abort()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "message": "invalid authorization header"})
            c.Abort()
            return
        }

        claims, err := jwt.ParseToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "message": "invalid token"})
            c.Abort()
            return
        }

        // 将用户信息存储到上下文中
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)

        c.Next()
    }
}

func AdminRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "message": "unauthorized"})
            c.Abort()
            return
        }

        if role.(int) != 1 { // 1 表示管理员角色
            c.JSON(http.StatusForbidden, gin.H{"code": 1003, "message": "permission denied"})
            c.Abort()
            return
        }

        c.Next()
    }
} 