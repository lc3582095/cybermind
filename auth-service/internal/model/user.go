package model

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID           int64          `gorm:"primaryKey" json:"id"`
    Username     string         `gorm:"size:50;not null;unique" json:"username"`
    Email        string         `gorm:"size:100;not null;unique" json:"email"`
    PasswordHash string         `gorm:"size:255;not null" json:"-"`
    AvatarURL    string         `gorm:"size:255" json:"avatar_url,omitempty"`
    Phone        string         `gorm:"size:20" json:"phone,omitempty"`
    Status       int            `gorm:"default:1" json:"status"`
    Role         int            `gorm:"default:0" json:"role"`
    Points       int            `gorm:"default:0" json:"points"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
    return "users"
} 