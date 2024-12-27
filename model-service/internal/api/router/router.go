package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"cybermind/model-service/internal/api/handler"
	"cybermind/model-service/internal/service"
)

// SetupRouter 设置路由
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 创建服务
	modelService := service.NewModelService(db)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 模型相关路由
		modelHandler := handler.NewModelHandler(modelService)
		v1.GET("/models", modelHandler.ListModels)
		v1.POST("/models", modelHandler.CreateModel)
		v1.PUT("/models/:id", modelHandler.UpdateModel)
		v1.PUT("/models/:id/status", modelHandler.UpdateModelStatus)
		v1.DELETE("/models/:id", modelHandler.DeleteModel)

		// 供应商相关路由
		providerHandler := handler.NewProviderHandler(db)
		v1.GET("/providers", providerHandler.ListProviders)
		v1.POST("/providers", providerHandler.CreateProvider)
		v1.PUT("/providers/:id", providerHandler.UpdateProvider)
		v1.PUT("/providers/:id/status", providerHandler.UpdateProviderStatus)
		v1.DELETE("/providers/:id", providerHandler.DeleteProvider)

		// API Key 池相关路由
		apiKeyPoolHandler := handler.NewAPIKeyPoolHandler(db)
		v1.GET("/api-keys", apiKeyPoolHandler.ListAPIKeyPools)
		v1.POST("/api-keys", apiKeyPoolHandler.CreateAPIKeyPool)
		v1.PUT("/api-keys/:id/status", apiKeyPoolHandler.UpdateAPIKeyPoolStatus)
		v1.DELETE("/api-keys/:id", apiKeyPoolHandler.DeleteAPIKeyPool)
	}

	return r
}
