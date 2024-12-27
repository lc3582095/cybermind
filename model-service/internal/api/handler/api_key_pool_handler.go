package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"cybermind/model-service/internal/model"
)

type APIKeyPoolHandler struct {
	db *gorm.DB
}

func NewAPIKeyPoolHandler(db *gorm.DB) *APIKeyPoolHandler {
	return &APIKeyPoolHandler{db: db}
}

// ListAPIKeyPools 获取 API Key 池列表
func (h *APIKeyPoolHandler) ListAPIKeyPools(c *gin.Context) {
	var keyPools []model.APIKeyPool
	if err := h.db.Find(&keyPools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1005,
			"message": "获取API Key池列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    keyPools,
	})
}

// CreateAPIKeyPool 创建 API Key
func (h *APIKeyPoolHandler) CreateAPIKeyPool(c *gin.Context) {
	var keyPool model.APIKeyPool
	if err := c.ShouldBindJSON(&keyPool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "参数错误",
		})
		return
	}

	// 设置初始状态为启用
	keyPool.Status = 1

	if err := h.db.Create(&keyPool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1005,
			"message": "创建API Key失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    keyPool,
	})
}

// UpdateAPIKeyPoolStatus 更新 API Key 状态
func (h *APIKeyPoolHandler) UpdateAPIKeyPoolStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "参数错误",
		})
		return
	}

	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "参数错误",
		})
		return
	}

	if err := h.db.Model(&model.APIKeyPool{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1005,
			"message": "更新API Key状态失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// DeleteAPIKeyPool 删除 API Key
func (h *APIKeyPoolHandler) DeleteAPIKeyPool(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "参数错误",
		})
		return
	}

	if err := h.db.Delete(&model.APIKeyPool{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1005,
			"message": "删除API Key失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
