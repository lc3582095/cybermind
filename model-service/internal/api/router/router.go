package router

import (
	"github.com/gin-gonic/gin"
	"cybermind/model-service/internal/api/handler"
)

// SetupRouter 配置路由
func SetupRouter(r *gin.Engine, modelHandler *handler.ModelHandler) {
	// API版本v1
	v1 := r.Group("/api/v1")
	{
		// 模型管理相关路由
		models := v1.Group("/models")
		{
			models.POST("", modelHandler.CreateModel)           // 创建模型
			models.GET("", modelHandler.ListModels)            // 获取模型列表
			models.GET("/:id", modelHandler.GetModel)          // 获取单个模型
			models.PUT("/:id", modelHandler.UpdateModel)       // 更新模型
			models.DELETE("/:id", modelHandler.DeleteModel)    // 删除模型
			models.PUT("/:id/status", modelHandler.UpdateModelStatus) // 更新模型状态
		}
	}
} 