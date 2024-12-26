package middleware

import (
	"net/http"
	"strings"

	"cybermind/admin-service/internal/model"
	"cybermind/admin-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Auth JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, model.Response{
				Code:    model.Unauthorized,
				Message: "未授权访问",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, model.Response{
				Code:    model.Unauthorized,
				Message: "无效的认证格式",
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.Response{
				Code:    model.Unauthorized,
				Message: "无效的Token",
			})
			c.Abort()
			return
		}

		// 将用户信息保存到上下文
		c.Set("admin_id", claims.ID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireRole 角色验证中间件
func RequireRole(role int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.Response{
				Code:    model.Unauthorized,
				Message: "未授权访问",
			})
			c.Abort()
			return
		}

		if userRole.(int) < role {
			c.JSON(http.StatusForbidden, model.Response{
				Code:    model.Forbidden,
				Message: "权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
} 