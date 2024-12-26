package model

import (
	"encoding/json"
	"time"
)

// Model 表示AI模型配置
type Model struct {
	ID               int64           `gorm:"primaryKey" json:"id"`
	Name             string          `gorm:"size:50;not null" json:"name"`
	Provider         string          `gorm:"size:50;not null" json:"provider"`
	APIType          string          `gorm:"size:20;not null;column:api_type" json:"api_type"`
	BaseURL          string          `gorm:"size:255;not null" json:"base_url"`
	APIKey           string          `gorm:"size:255;not null" json:"api_key"`
	ModelName        string          `gorm:"size:50;not null" json:"model_name"`
	PointsPerRequest int             `gorm:"not null;column:points_per_request" json:"points_per_request"`
	Config           json.RawMessage `gorm:"type:jsonb" json:"config"`
	Status           int             `gorm:"default:1" json:"status"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// ModelConfig 表示模型配置参数
type ModelConfig struct {
	Temperature      float64 `json:"temperature,omitempty"`
	MaxTokens        int     `json:"max_tokens,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
}

// TableName 指定表名
func (Model) TableName() string {
	return "models"
}

// DefaultModels 返回默认的模型配置列表
func DefaultModels() []Model {
	defaultConfig := DefaultModelConfigs[ProviderOpenAI]
	configMap := map[string]interface{}{
		"temperature":       defaultConfig.Temperature,
		"max_tokens":       defaultConfig.MaxTokens,
		"top_p":            defaultConfig.TopP,
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