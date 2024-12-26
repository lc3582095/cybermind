package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"cybermind/model-service/internal/model"
	"cybermind/model-service/internal/service"
)

type ModelHandler struct {
	modelService *service.ModelService
}

func NewModelHandler(modelService *service.ModelService) *ModelHandler {
	return &ModelHandler{modelService: modelService}
}

// CreateModel 创建模型配置
func (h *ModelHandler) CreateModel(c *gin.Context) {
	var m model.Model
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1001, "message": "Invalid parameters", "error": err.Error()})
		return
	}

	if err := h.modelService.CreateModel(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1005, "message": "Failed to create model", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": m})
}

// GetModel 获取模型配置
func (h *ModelHandler) GetModel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1001, "message": "Invalid ID", "error": err.Error()})
		return
	}

	m, err := h.modelService.GetModel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1004, "message": "Model not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": m})
}

// UpdateModel 更新模型配置
func (h *ModelHandler) UpdateModel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1001, "message": "Invalid ID", "error": err.Error()})
		return
	}

	var m model.Model
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1001, "message": "Invalid parameters", "error": err.Error()})
		return
	}
	m.ID = id

	if err := h.modelService.UpdateModel(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1005, "message": "Failed to update model", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": m})
}

// ListModels 获取模型列表
func (h *ModelHandler) ListModels(c *gin.Context) {
	models, err := h.modelService.ListModels(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1005, "message": "Failed to get models", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": models})
}

// DeleteModel 删除模型
func (h *ModelHandler) DeleteModel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1001, "message": "Invalid ID", "error": err.Error()})
		return
	}

	if err := h.modelService.DeleteModel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1005, "message": "Failed to delete model", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// UpdateModelStatus 更新模型状态
func (h *ModelHandler) UpdateModelStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1001, "message": "Invalid ID", "error": err.Error()})
		return
	}

	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1001, "message": "Invalid parameters", "error": err.Error()})
		return
	}

	if err := h.modelService.UpdateModelStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1005, "message": "Failed to update model status", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
} 