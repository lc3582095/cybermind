package router

import (
	"github.com/gin-gonic/gin"
	"cybermind/chat-service/internal/api/handler"
	"cybermind/chat-service/pkg/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	// 创建处理器实例
	chatHandler := handler.NewChatHandler()

	// API路由组
	api := r.Group("/api/v1")
	{
		// 对话相关路由(需要认证)
		conversations := api.Group("/conversations", middleware.AuthMiddleware())
		{
			conversations.POST("", chatHandler.CreateConversation)           // 创建对话
			conversations.GET("", chatHandler.ListConversations)            // 获取对话列表
			conversations.GET("/detail/:id", chatHandler.GetConversation)   // 获取对话详情
			conversations.GET("/messages/:conversation_id", chatHandler.GetMessages)  // 获取消息列表
			conversations.POST("/messages", chatHandler.AddMessage)         // 添加消息
		}
	}
} 