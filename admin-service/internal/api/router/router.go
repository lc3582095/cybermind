package router

import (
	"cybermind/admin-service/internal/api/handler"
	"cybermind/admin-service/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handler.Login)
			auth.POST("/logout", middleware.Auth(), handler.Logout)
			auth.GET("/info", middleware.Auth(), handler.GetAdminInfo)
			auth.PUT("/password", middleware.Auth(), handler.UpdatePassword)
		}

		// 管理员相关路由
		admin := v1.Group("/admin", middleware.Auth())
		{
			// 管理员管理
			admin.GET("/list", middleware.RequireRole(2), handler.GetAdminList)
			admin.POST("/create", middleware.RequireRole(2), handler.CreateAdmin)
			admin.PUT("/:id", middleware.RequireRole(2), handler.UpdateAdmin)
			admin.DELETE("/:id", middleware.RequireRole(2), handler.DeleteAdmin)

			// 用户管理
			admin.GET("/users", handler.GetUserList)
			admin.PUT("/users/:id/status", handler.UpdateUserStatus)
			admin.GET("/users/:id/detail", handler.GetUserDetail)

			// 模型管理
			admin.GET("/models", handler.GetModelList)
			admin.POST("/models", middleware.RequireRole(2), handler.CreateModel)
			admin.PUT("/models/:id", middleware.RequireRole(2), handler.UpdateModel)
			admin.DELETE("/models/:id", middleware.RequireRole(2), handler.DeleteModel)

			// 订单管理
			admin.GET("/orders", handler.GetOrderList)
			admin.GET("/orders/:id", handler.GetOrderDetail)
			admin.PUT("/orders/:id/status", handler.UpdateOrderStatus)

			// 支付管理
			admin.GET("/payments", handler.GetPaymentList)
			admin.GET("/payments/:id", handler.GetPaymentDetail)
			admin.POST("/payments/:id/refund", middleware.RequireRole(2), handler.CreateRefund)

			// 统计分析
			admin.GET("/stats/overview", handler.GetStatsOverview)
			admin.GET("/stats/daily", handler.GetDailyStats)
			admin.GET("/stats/user", handler.GetUserStats)
			admin.GET("/stats/order", handler.GetOrderStats)
		}
	}

	return r
} 