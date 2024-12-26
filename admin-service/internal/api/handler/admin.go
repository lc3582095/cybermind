package handler

import (
	"net/http"
	"strconv"

	"cybermind/admin-service/internal/model"
	"cybermind/admin-service/pkg/database"
	"cybermind/admin-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateAdminRequest 创建管理员请求
type CreateAdminRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     int    `json:"role" binding:"required,oneof=1 2"`
}

// CreateAdmin 创建管理员
func CreateAdmin(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	// 检查邮箱是否已存在
	var count int64
	database.DB.Model(&model.Admin{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "邮箱已存在",
		})
		return
	}

	// 密码加密
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	admin := &model.Admin{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         req.Role,
		Status:       1,
	}

	if err := database.DB.Create(admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "创建成功",
		Data:    admin,
	})
}

// GetAdminList 获取管理员列表
func GetAdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	username := c.Query("username")

	var admins []model.Admin
	var total int64
	query := database.DB.Model(&model.Admin{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	query.Count(&total)
	if err := query.Offset((page - 1) * size).Limit(size).Find(&admins).Error; err != nil {
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
			List:  admins,
		},
	})
}

// UpdateAdminRequest 更新管理员请求
type UpdateAdminRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Role     int    `json:"role" binding:"required,oneof=1 2"`
	Status   int    `json:"status" binding:"required,oneof=0 1"`
}

// UpdateAdmin 更新管理员
func UpdateAdmin(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	var admin model.Admin
	if err := database.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{
			Code:    model.NotFound,
			Message: "管理员不存在",
		})
		return
	}

	admin.Username = req.Username
	admin.Role = req.Role
	admin.Status = req.Status

	if err := database.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "更新成功",
		Data:    admin,
	})
}

// DeleteAdmin 删除管理员
func DeleteAdmin(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	var admin model.Admin
	if err := database.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{
			Code:    model.NotFound,
			Message: "管理员不存在",
		})
		return
	}

	if err := database.DB.Delete(&admin).Error; err != nil {
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