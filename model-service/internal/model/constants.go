package model

// 模型提供商
const (
	ProviderOpenAI = "OpenAI"
)

// API基础URL
const (
	OpenAIBaseURL = "https://api.fast-tunnel.one"
)

// 默认模型配置
var DefaultModelConfigs = map[string]ModelConfig{
	ProviderOpenAI: {
		Temperature:      0.7,
		MaxTokens:       2000,
		TopP:            1.0,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
	},
}

// API密钥
const (
	OpenAIKey = "sk-kUNLgEMdDBf1LHWLB7C90674A2A54b99BeF5D5E8D39fD04e"
) 