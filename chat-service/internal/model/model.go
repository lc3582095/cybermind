package model

import (
	"time"
	"gorm.io/gorm"
)

// Conversation 对话模型
type Conversation struct {
	ID             int64          `gorm:"primaryKey" json:"id"`
	UserID         int64          `gorm:"not null;index" json:"user_id"`
	ModelID        int64          `gorm:"not null" json:"model_id"`
	Title          string         `gorm:"size:255" json:"title"`
	PointsConsumed int            `gorm:"not null;default:0" json:"points_consumed"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Messages       []Message      `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

// Message 消息模型
type Message struct {
	ID             int64          `gorm:"primaryKey" json:"id"`
	ConversationID int64          `gorm:"not null;index" json:"conversation_id"`
	Role           string         `gorm:"size:20;not null" json:"role"` // system/user/assistant
	Content        string         `gorm:"type:text;not null" json:"content"`
	TokensCount    int            `json:"tokens_count,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
} 