# Part 2: AI模型管理模块开发文档

## 1. 模块概述

AI模型管理模块主要负责管理和配置不同AI服务提供商的模型信息，包括API密钥、配置参数等。该模块提供了完整的REST API接口，支持模型的增删改查等基本操作，并实现了API密钥池和供应商管理功能。

## 2. 技术栈

- 后端框架：Gin
- 数据库：PostgreSQL
- ORM：GORM
- 配置格式：JSONB
- 安全：JWT认证、HTTPS、API网关

## 3. 目录结构

```
cybermind/model-service/
├── cmd/
│   └── server/
│       └── main.go           # 程序入口
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   ├── model_handler.go    # 模型处理器
│   │   │   ├── provider_handler.go  # 供应商处理器
│   │   │   └── api_key_pool_handler.go  # API密钥池处理器
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
    ProxyURL         string          `gorm:"size:255" json:"proxy_url"`
    APIKey           string          `gorm:"size:255;not null" json:"api_key"`
    ModelName        string          `gorm:"size:50;not null" json:"model_name"`
    PointsPerRequest int             `gorm:"not null;column:points_per_request" json:"points_per_request"`
    Tags             pq.StringArray  `gorm:"type:text[];column:tags" json:"tags"`
    Config           json.RawMessage `gorm:"type:jsonb" json:"config"`
    Preset           string          `gorm:"type:text" json:"preset"`
    Status           int             `gorm:"default:1" json:"status"`
    CreatedAt        time.Time       `json:"created_at"`
    UpdatedAt        time.Time       `json:"updated_at"`
    DeletedAt        gorm.DeletedAt  `gorm:"index" json:"-"`
}
```

### Provider 结构体
```go
type Provider struct {
    ID         int64          `gorm:"primaryKey" json:"id"`
    Name       string         `gorm:"size:50;not null" json:"name"`
    Code       string         `gorm:"size:50;not null;unique" json:"code"`
    ModelCount int            `gorm:"default:0" json:"model_count"`
    SortID     int            `gorm:"default:0" json:"sort_id"`
    Status     int            `gorm:"default:1" json:"status"`
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### APIKeyPool 结构体
```go
type APIKeyPool struct {
    ID         int64          `gorm:"primaryKey" json:"id"`
    ModelID    int64          `gorm:"not null;index" json:"model_id"`
    APIKey     string         `gorm:"size:255;not null" json:"api_key"`
    Status     int            `gorm:"default:1" json:"status"`
    UsageCount int64          `gorm:"default:0" json:"usage_count"`
    LastUsedAt *time.Time     `json:"last_used_at,omitempty"`
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## 5. API接口

### 5.1 模型管理接口

#### 获取模型列表
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
            "base_url": "https://api.openai.com/v1",
            "proxy_url": "https://api.fast-tunnel.one",
            "api_key": "sk-xxx",
            "model_name": "gpt-4",
            "points_per_request": 10,
            "tags": ["对话", "文本"],
            "config": {
                "temperature": 0.7,
                "max_tokens": 2000,
                "top_p": 1
            },
            "preset": "你是一个AI助手...",
            "status": 1,
            "created_at": "2024-12-24T10:12:02.807826Z",
            "updated_at": "2024-12-24T10:12:02.807826Z"
        }
    ]
}
```

### 5.2 供应商管理接口

#### 获取供应商列表
- 路径: GET `/api/v1/providers`
- 响应示例:
```json
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "id": 1,
            "name": "OpenAI",
            "code": "openai",
            "model_count": 5,
            "sort_id": 1000,
            "status": 1,
            "created_at": "2024-12-24T10:12:02.807826Z",
            "updated_at": "2024-12-24T10:12:02.807826Z"
        }
    ]
}
```

### 5.3 API密钥池接口

#### 获取API密钥列表
- 路径: GET `/api/v1/api-keys`
- 响应示例:
```json
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "id": 1,
            "model_id": 1,
            "api_key": "sk-xxx",
            "status": 1,
            "usage_count": 100,
            "last_used_at": "2024-12-24T10:12:02.807826Z",
            "created_at": "2024-12-24T10:12:02.807826Z",
            "updated_at": "2024-12-24T10:12:02.807826Z"
        }
    ]
}
```

## 6. 安全措施

### 6.1 API安全
- 使用 HTTPS 加密传输
- 实现 JWT 身份认证
- 添加请求频率限制
- 实现 API 调用审计日志

### 6.2 数据安全
- API密钥等敏感信息加密存储
- 数据库访问权限控制
- 定期数据备份
- 软删除机制

### 6.3 部署安全
- 使用反向代理隐藏真实服务
- 配置合适的CORS策略
- 实现API网关层
- 监控异常访问行为

## 7. 错误处理

### 错误响应格式
```json
{
    "code": 1001,
    "message": "Invalid parameters",
    "error": "具体错误信息"
}
```

### 错误码说明
- 0: 成功
- 1001: 参数无效
- 1002: 未授权
- 1003: 禁止访问
- 1004: 资源不存在
- 1005: 服务器内部错误

## 8. 部署说明

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
- 运行模式: debug/release
- 日志级别: info/debug/error

### 8.3 反向代理配置
```nginx
location /api/ {
    proxy_pass http://localhost:8081/;
    proxy_hide_header Server;
    proxy_hide_header X-Powered-By;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

## 9. 开发注意事项

1. Config字段使用JSONB类型存储，支持灵活的配置参数
2. API密钥等敏感信息在实际部署时应该使用环境变量或配置中心
3. 所有API都实现了统一的错误处理和日志记录
4. 数据库操作使用GORM事务确保数据一致性
5. 服务启动时会自动进行数据库迁移和默认数据初始化
6. 实现了API密钥池的轮询调用机制
7. 支持模型预设和标签管理
8. 提供了完整的供应商管理功能 