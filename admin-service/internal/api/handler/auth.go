package handler

import (
	"net/http"

	"cybermind/admin-service/internal/model"
	"cybermind/admin-service/pkg/database"
	"cybermind/admin-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login 管理员登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	var admin model.Admin
	if err := database.DB.Where("email = ?", req.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.AdminNotExist,
			Message: "管理员不存在",
		})
		return
	}

	if !utils.CheckPassword(req.Password, admin.PasswordHash) {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.WrongPassword,
			Message: "密码错误",
		})
		return
	}

	if admin.Status == 0 {
		c.JSON(http.StatusForbidden, model.Response{
			Code:    model.AdminDisabled,
			Message: "账号已禁用",
		})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(admin.ID, admin.Username, admin.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "登录成功",
		Data: model.LoginResponse{
			Token: token,
			Admin: admin,
		},
	})
}

// Logout 管理员登出
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "登出成功",
	})
}

// GetAdminInfo 获取管理员信息
func GetAdminInfo(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var admin model.Admin
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{
			Code:    model.NotFound,
			Message: "管理员不存在",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "获取成功",
		Data:    admin,
	})
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UpdatePassword 更新密码
func UpdatePassword(c *gin.Context) {
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	adminID, _ := c.Get("admin_id")

	var admin model.Admin
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{
			Code:    model.NotFound,
			Message: "管理员不存在",
		})
		return
	}

	if !utils.CheckPassword(req.OldPassword, admin.PasswordHash) {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.WrongPassword,
			Message: "原密码错误",
		})
		return
	}

	// 更新密码
	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	admin.PasswordHash = newHash
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
	})
} 