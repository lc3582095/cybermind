package handler

import (
	"net/http"
	"strconv"

	"cybermind/admin-service/internal/model"
	"cybermind/admin-service/pkg/database"

	"github.com/gin-gonic/gin"
)

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	username := c.Query("username")
	email := c.Query("email")
	status := c.Query("status")

	var users []struct {
		ID        int64  `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Status    int    `json:"status"`
		Points    int    `json:"points"`
		CreatedAt string `json:"created_at"`
	}
	var total int64
	query := database.DB.Table("users")

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	if err := query.Offset((page - 1) * size).Limit(size).
		Select("id, username, email, phone, status, points, created_at").
		Find(&users).Error; err != nil {
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
			List:  users,
		},
	})
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	Status int    `json:"status" binding:"required,oneof=0 1"`
	Reason string `json:"reason" binding:"required"`
}

// UpdateUserStatus 更新用户状态
func UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	var req UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	// 更新用户状态
	if err := database.DB.Table("users").Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": req.Status,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	// 记录操作日志
	adminID, _ := c.Get("admin_id")
	operation := &model.AdminOperation{
		AdminID:     adminID.(int64),
		Module:      "user",
		Action:      "update_status",
		Description: req.Reason,
		IP:          c.ClientIP(),
	}
	database.DB.Create(operation)

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "更新成功",
	})
}

// GetUserDetail 获取用户详情
func GetUserDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	// 用户基本信息
	var user struct {
		ID        int64  `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Status    int    `json:"status"`
		Points    int    `json:"points"`
		CreatedAt string `json:"created_at"`
	}
	if err := database.DB.Table("users").
		Select("id, username, email, phone, status, points, created_at").
		Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{
			Code:    model.NotFound,
			Message: "用户不存在",
		})
		return
	}

	// 用户套餐信息
	var packages []struct {
		PackageName      string `json:"package_name"`
		PointsRemaining  int    `json:"points_remaining"`
		StartTime       string `json:"start_time"`
		EndTime         string `json:"end_time"`
		Status          int    `json:"status"`
	}
	database.DB.Table("user_packages").
		Select("packages.name as package_name, user_packages.points_remaining, user_packages.start_time, user_packages.end_time, user_packages.status").
		Joins("left join packages on packages.id = user_packages.package_id").
		Where("user_packages.user_id = ?", id).
		Find(&packages)

	// 最近的对话记录
	var conversations []struct {
		ID             int64  `json:"id"`
		Title          string `json:"title"`
		ModelName      string `json:"model_name"`
		PointsConsumed int    `json:"points_consumed"`
		CreatedAt      string `json:"created_at"`
	}
	database.DB.Table("conversations").
		Select("conversations.id, conversations.title, models.name as model_name, conversations.points_consumed, conversations.created_at").
		Joins("left join models on models.id = conversations.model_id").
		Where("conversations.user_id = ?", id).
		Order("conversations.created_at desc").
		Limit(10).
		Find(&conversations)

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "获取成功",
		Data: gin.H{
			"user":          user,
			"packages":      packages,
			"conversations": conversations,
		},
	})
} 