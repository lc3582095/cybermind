package handler

import (
	"net/http"
	"strconv"
	"time"

	"cybermind/admin-service/internal/model"
	"cybermind/admin-service/pkg/database"

	"github.com/gin-gonic/gin"
)

// CreateModelRequest 创建模型请求
type CreateModelRequest struct {
	Name             string `json:"name" binding:"required"`
	Provider         string `json:"provider" binding:"required"`
	APIType          string `json:"api_type" binding:"required"`
	BaseURL          string `json:"base_url" binding:"required"`
	APIKey           string `json:"api_key" binding:"required"`
	ModelName        string `json:"model_name" binding:"required"`
	PointsPerRequest int    `json:"points_per_request" binding:"required,min=1"`
}

// CreateModel 创建模型
func CreateModel(c *gin.Context) {
	var req CreateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	// 检查模型名称是否已存在
	var count int64
	database.DB.Table("models").Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "模型名称已存在",
		})
		return
	}

	// 创建模型
	modelData := map[string]interface{}{
		"name":               req.Name,
		"provider":           req.Provider,
		"api_type":          req.APIType,
		"base_url":          req.BaseURL,
		"api_key":           req.APIKey,
		"model_name":        req.ModelName,
		"points_per_request": req.PointsPerRequest,
		"status":            1,
		"created_at":        time.Now(),
	}

	if err := database.DB.Table("models").Create(modelData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "创建成功",
		Data:    modelData,
	})
}

// GetModelList 获取模型列表
func GetModelList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	name := c.Query("name")
	provider := c.Query("provider")
	status := c.Query("status")

	var models []struct {
		ID                int64     `json:"id"`
		Name              string    `json:"name"`
		Provider          string    `json:"provider"`
		APIType           string    `json:"api_type"`
		BaseURL           string    `json:"base_url"`
		ModelName         string    `json:"model_name"`
		PointsPerRequest  int       `json:"points_per_request"`
		Status            int       `json:"status"`
		CreatedAt         time.Time `json:"created_at"`
		LastRequestTime   time.Time `json:"last_request_time"`
		TotalRequests     int64     `json:"total_requests"`
		SuccessRequests   int64     `json:"success_requests"`
		FailedRequests    int64     `json:"failed_requests"`
		AverageLatency    float64   `json:"average_latency"`
		TotalTokens       int64     `json:"total_tokens"`
		TotalPoints       int64     `json:"total_points"`
	}

	var total int64
	query := database.DB.Table("models")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if provider != "" {
		query = query.Where("provider = ?", provider)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	if err := query.Offset((page - 1) * size).Limit(size).Find(&models).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "获取成功",
		Data: model.PageResponse{
			Total: total,
			List:  models,
		},
	})
}

// UpdateModelRequest 更新模型请求
type UpdateModelRequest struct {
	Name             string `json:"name" binding:"required"`
	BaseURL          string `json:"base_url" binding:"required"`
	APIKey           string `json:"api_key" binding:"required"`
	PointsPerRequest int    `json:"points_per_request" binding:"required,min=1"`
	Status           int    `json:"status" binding:"required,oneof=0 1"`
}

// UpdateModel 更新模型
func UpdateModel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	var req UpdateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	// 检查模型是否存在
	var count int64
	database.DB.Table("models").Where("id = ?", id).Count(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, model.Response{
			Code:    model.NotFound,
			Message: "模型不存在",
		})
		return
	}

	// 检查名称是否重复
	database.DB.Table("models").Where("name = ? AND id != ?", req.Name, id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "模型名称已存在",
		})
		return
	}

	// 更新模型
	modelData := map[string]interface{}{
		"name":               req.Name,
		"base_url":          req.BaseURL,
		"api_key":           req.APIKey,
		"points_per_request": req.PointsPerRequest,
		"status":            req.Status,
		"updated_at":        time.Now(),
	}

	if err := database.DB.Table("models").Where("id = ?", id).Updates(modelData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "更新成功",
	})
}

// DeleteModel 删除模型
func DeleteModel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	// 检查模型是否存在
	var count int64
	database.DB.Table("models").Where("id = ?", id).Count(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, model.Response{
			Code:    model.NotFound,
			Message: "模型不存在",
		})
		return
	}

	// 删除模型
	if err := database.DB.Table("models").Where("id = ?", id).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "删除成功",
	})
} 