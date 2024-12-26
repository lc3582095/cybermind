package model

import (
	"time"

	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	ID           int64          `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:50;not null;unique" json:"username"`
	Email        string         `gorm:"size:100;not null;unique" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Role         int            `gorm:"default:1" json:"role"` // 1: 普通管理员, 2: 超级管理员
	Status       int            `gorm:"default:1" json:"status"` // 0: 禁用, 1: 正常
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// AdminLoginHistory 管理员登录历史
type AdminLoginHistory struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	AdminID   int64     `gorm:"not null" json:"admin_id"`
	IP        string    `gorm:"size:50" json:"ip"`
	UserAgent string    `gorm:"size:255" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// AdminOperation 管理员操作日志
type AdminOperation struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	AdminID     int64     `gorm:"not null" json:"admin_id"`
	Module      string    `gorm:"size:50" json:"module"`      // 操作模块
	Action      string    `gorm:"size:50" json:"action"`      // 操作类型
	Description string    `gorm:"size:255" json:"description"` // 操作描述
	IP          string    `gorm:"size:50" json:"ip"`          // 操作IP
	CreatedAt   time.Time `json:"created_at"`
} 