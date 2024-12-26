package handler

import (
    "net/http"
    "cybermind/auth-service/internal/service"
    "cybermind/auth-service/pkg/jwt"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

// Register 用户注册处理
func (h *UserHandler) Register(c *gin.Context) {
    var req service.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    1001,
            "message": "invalid request parameters",
        })
        return
    }

    user, err := h.userService.Register(&req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    1001,
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    0,
        "message": "success",
        "data": gin.H{
            "user_id":  user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}

// Login 用户登录处理
func (h *UserHandler) Login(c *gin.Context) {
    var req service.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    1001,
            "message": "invalid request parameters",
        })
        return
    }

    user, err := h.userService.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "code":    1002,
            "message": err.Error(),
        })
        return
    }

    // 生成JWT token
    token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":    1005,
            "message": "failed to generate token",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    0,
        "message": "success",
        "data": gin.H{
            "token": token,
            "user": gin.H{
                "id":       user.ID,
                "username": user.Username,
                "email":    user.Email,
                "points":   user.Points,
            },
        },
    })
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
    userID, _ := c.Get("userID")
    user, err := h.userService.GetUserByID(userID.(int64))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "code":    1004,
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    0,
        "message": "success",
        "data": gin.H{
            "id":         user.ID,
            "username":   user.Username,
            "email":      user.Email,
            "phone":      user.Phone,
            "points":     user.Points,
            "avatar_url": user.AvatarURL,
        },
    })
} 