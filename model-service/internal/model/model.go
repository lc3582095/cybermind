package model

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Model 表示AI模型配置
type Model struct {
	ID               int64           `gorm:"primaryKey" json:"id"`
	Name             string          `gorm:"size:50;not null" json:"name"`
	Provider         string          `gorm:"size:50;not null" json:"provider"`
	APIType          string          `gorm:"size:20;not null;column:api_type" json:"api_type"`
	BaseURL          string          `gorm:"size:255;not null" json:"base_url"`
	ProxyURL         string          `gorm:"size:255" json:"proxy_url"`
	APIKey           string          `gorm:"size:255;not null" json:"api_key"`
	ModelName        string          `gorm:"size:50;not null" json:"model_name"`
	PointsPerRequest int             `gorm:"not null;column:points_per_request" json:"points_per_request"`
	Tags             pq.StringArray  `gorm:"type:text[];column:tags" json:"tags"`
	Config           json.RawMessage `gorm:"type:jsonb" json:"config"`
	Preset           string          `gorm:"type:text" json:"preset"` // 模型预设描述
	Status           int             `gorm:"default:1" json:"status"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"-"`
}

// ModelConfig 表示模型配置参数
type ModelConfig struct {
	Temperature      float64 `json:"temperature,omitempty"`
	MaxTokens        int     `json:"max_tokens,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
}

// Provider 供应商模型
type Provider struct {
	ID         int64          `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:50;not null" json:"name"`        // 供应商名称
	Code       string         `gorm:"size:50;not null;unique" json:"code"` // 供应商代码
	ModelCount int            `gorm:"default:0" json:"model_count"`        // 模型数量
	SortID     int            `gorm:"default:0" json:"sort_id"`            // 排序ID
	Status     int            `gorm:"default:1" json:"status"`             // 状态：1-启用，0-禁用
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// APIKeyPool API密钥池配置
type APIKeyPool struct {
	ID         int64          `gorm:"primaryKey" json:"id"`
	ModelID    int64          `gorm:"not null;index" json:"model_id"`
	APIKey     string         `gorm:"size:255;not null" json:"api_key"`
	Status     int            `gorm:"default:1" json:"status"`      // 状态：1-启用，0-禁用
	UsageCount int64          `gorm:"default:0" json:"usage_count"` // 使用次数
	LastUsedAt *time.Time     `json:"last_used_at,omitempty"`       // 最后使用时间
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Model) TableName() string {
	return "models"
}

// TableName 指定表名
func (Provider) TableName() string {
	return "providers"
}

// TableName 指定表名
func (APIKeyPool) TableName() string {
	return "api_key_pools"
}

// DefaultModels 返回默认的模型配置列表
func DefaultModels() []Model {
	defaultConfig := DefaultModelConfigs[ProviderOpenAI]
	configMap := map[string]interface{}{
		"temperature":       defaultConfig.Temperature,
		"max_tokens":        defaultConfig.MaxTokens,
		"top_p":             defaultConfig.TopP,
		"frequency_penalty": defaultConfig.FrequencyPenalty,
		"presence_penalty":  defaultConfig.PresencePenalty,
	}

	configJSON, _ := json.Marshal(configMap)

	return []Model{
		{
			Name:             "GPT-4",
			Provider:         ProviderOpenAI,
			APIType:          "chat/completions",
			BaseURL:          OpenAIBaseURL,
			APIKey:           OpenAIKey,
			ModelName:        "gpt-4",
			PointsPerRequest: 10,
			Config:           configJSON,
			Status:           1,
		},
	}
}

// UpdateModelCount 更新供应商的模型数量
func (p *Provider) UpdateModelCount(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Model{}).Where("provider = ?", p.Code).Count(&count).Error; err != nil {
		return err
	}
	p.ModelCount = int(count)
	return db.Model(p).Update("model_count", p.ModelCount).Error
}
