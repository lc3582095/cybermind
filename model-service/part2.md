# Part 2: AI模型管理模块开发文档

## 1. 模块概述

AI模型管理模块主要负责管理和配置不同AI服务提供商的模型信息，包括API密钥、配置参数等。该模块提供了完整的REST API接口，支持模型的增删改查等基本操作。

## 2. 技术栈

- 后端框架：Gin
- 数据库：PostgreSQL
- ORM：GORM
- 配置格式：JSONB

## 3. 目录结构

```
cybermind/model-service/
├── cmd/
│   └── server/
│       └── main.go           # 程序入口
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   └── model_handler.go  # API处理器
│   │   └── router/
│   │       └── router.go     # 路由配置
│   ├── model/
│   │   ├── model.go         # 数据模型定义
│   │   └── constants.go     # 常量定义
│   └── service/
│       └── model_service.go  # 业务逻辑
└── pkg/
    ├── database/            # 数据库相关
    └── middleware/         # 中间件
```

## 4. 数据模型

### Model 结构体
```go
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
```

### ModelConfig 结构体
```go
type ModelConfig struct {
    Temperature      float64 `json:"temperature,omitempty"`
    MaxTokens        int     `json:"max_tokens,omitempty"`
    TopP             float64 `json:"top_p,omitempty"`
    FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
    PresencePenalty  float64 `json:"presence_penalty,omitempty"`
}
```

## 5. API接口

### 5.1 获取模型列表
- 路径: GET `/api/v1/models`
- 响应示例:
```json
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "id": 1,
            "name": "GPT-4",
            "provider": "OpenAI",
            "api_type": "chat/completions",
            "base_url": "https://api.fast-tunnel.one",
            "api_key": "sk-xxx",
            "model_name": "gpt-4",
            "points_per_request": 10,
            "config": {
                "temperature": 0.7,
                "max_tokens": 2000,
                "top_p": 1
            },
            "status": 1,
            "created_at": "2024-12-24T10:12:02.807826Z",
            "updated_at": "2024-12-24T10:12:02.807826Z"
        }
    ]
}
```

### 5.2 创建模型
- 路径: POST `/api/v1/models`
- 请求示例:
```json
{
    "name": "GPT-4-Turbo",
    "provider": "OpenAI",
    "api_type": "chat/completions",
    "base_url": "https://api.fast-tunnel.one",
    "api_key": "sk-xxx",
    "model_name": "gpt-4-turbo",
    "points_per_request": 12,
    "config": {
        "temperature": 0.8,
        "max_tokens": 4000,
        "top_p": 1
    },
    "status": 1
}
```

### 5.3 更新模型状态
- 路径: PUT `/api/v1/models/:id/status`
- 请求示例:
```json
{
    "status": 0
}
```

### 5.4 删除模型
- 路径: DELETE `/api/v1/models/:id`
- 响应示例:
```json
{
    "code": 0,
    "message": "success"
}
```

## 6. 测试用例

### 6.1 获取模型列表测试
```bash
curl -v http://localhost:8081/api/v1/models
```

### 6.2 创建模型测试
```bash
curl -X POST http://localhost:8081/api/v1/models \
  -H "Content-Type: application/json" \
  -d '{
    "name": "GPT-4-Turbo",
    "provider": "OpenAI",
    "api_type": "chat/completions",
    "base_url": "https://api.fast-tunnel.one",
    "api_key": "sk-kUNLgEMdDBf1LHWLB7C90674A2A54b99BeF5D5E8D39fD04e",
    "model_name": "gpt-4-turbo",
    "points_per_request": 12,
    "config": {
        "temperature": 0.8,
        "max_tokens": 4000,
        "top_p": 1
    },
    "status": 1
}'
```

### 6.3 更新模型状态测试
```bash
curl -X PUT http://localhost:8081/api/v1/models/4/status \
  -H "Content-Type: application/json" \
  -d '{"status": 0}'
```

### 6.4 删除模型测试
```bash
curl -X DELETE http://localhost:8081/api/v1/models/4
```

## 7. 错误处理

### 错误响应格式
```json
{
    "code": 1001,
    "message": "Invalid parameters",
    "error": "具体错误信息"
}
```

### 错��码说明
- 0: 成功
- 1001: 参数无效
- 1004: 资源不存在
- 1005: 服务器内部错误

## 8. 配置说明

### 8.1 数据库配置
```
host: cybermind-postgresql.ns-han88ija.svc
user: postgres
password: wkzhx7jn
dbname: postgres
port: 5432
```

### 8.2 服务配置
- 监听端口: 8081
- 运行模式: debug

## 9. 开发注意事项

1. Config字段使用JSONB类型存储，支持灵活的配置参数
2. API密钥等敏感信息在实际部署时应该使用环境变量或配置中心
3. 所有API都实现了统一的错误处理和日志记录
4. 数据库操作使用GORM事务确保数据一致性
5. 服务启动时会自动进行数据库迁移和默认数据初始化 