package middleware

import (
	"bytes"
	"io"
	"time"

	"cybermind/admin-service/internal/model"
	"cybermind/admin-service/pkg/database"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 处理请求
		c.Next()

		// 如果是登录成功的请求，记录登录历史
		if c.FullPath() == "/api/v1/auth/login" && c.Writer.Status() == 200 {
			adminID, exists := c.Get("admin_id")
			if exists {
				history := &model.AdminLoginHistory{
					AdminID:   adminID.(int64),
					IP:        c.ClientIP(),
					UserAgent: c.Request.UserAgent(),
				}
				database.DB.Create(history)
			}
		}

		// 记录操作日志
		if adminID, exists := c.Get("admin_id"); exists {
			operation := &model.AdminOperation{
				AdminID:     adminID.(int64),
				Module:     c.FullPath(),
				Action:     c.Request.Method,
				Description: string(requestBody),
				IP:         c.ClientIP(),
			}
			database.DB.Create(operation)
		}
	}
} 