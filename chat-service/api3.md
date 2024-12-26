# CyberMind对话管理服务API文档

## 1. 服务概述

### 1.1 基本信息
- 服务名称：对话管理服务（Chat Service）
- 基础URL：`http://localhost:8081/api/v1`
- 开发语言：Go
- 主要依赖：
  - gin-gonic/gin: Web框架
  - gorm.io/gorm: ORM框架
  - golang-jwt/jwt: JWT认证

### 1.2 数据库配置
- 数据库类型：PostgreSQL
- 连接信息：
  ```
  Host: dbconn.sealosbja.site
  Port: 37550
  User: postgres
  Password: wkzhx7jn
  Database: postgres
  ```

## 2. 数据模型

### 2.1 对话模型（Conversation）
```go
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
```

### 2.2 消息模型（Message）
```go
type Message struct {
    ID             int64          `gorm:"primaryKey" json:"id"`
    ConversationID int64          `gorm:"not null;index" json:"conversation_id"`
    Role           string         `gorm:"size:20;not null" json:"role"` // system/user/assistant
    Content        string         `gorm:"type:text;not null" json:"content"`
    TokensCount    int            `json:"tokens_count,omitempty"`
    CreatedAt      time.Time      `json:"created_at"`
    DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## 3. API接口

### 3.1 创建对话
- **接口**：`POST /conversations`
- **描述**：创建新的对话
- **请求头**：
  ```
  Content-Type: application/json
  Authorization: Bearer <token>
  ```
- **请求体**：
  ```json
  {
    "model_id": 1,        // AI模型ID
    "title": "测试对话"    // 对话标题
  }
  ```
- **响应**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "id": 1,
      "user_id": 1,
      "model_id": 1,
      "title": "测试对话",
      "points_consumed": 0,
      "created_at": "2024-12-24T11:55:36Z",
      "updated_at": "2024-12-24T11:55:36Z"
    }
  }
  ```

### 3.2 获取对话列表
- **接口**：`GET /conversations`
- **描述**：获取用户的对话列表
- **请求头**：
  ```
  Authorization: Bearer <token>
  ```
- **查询参数**：
  - page: 页码（默认1）
  - size: 每页数量（默认20）
- **响应**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "total": 100,
      "list": [
        {
          "id": 1,
          "user_id": 1,
          "model_id": 1,
          "title": "测试对话",
          "points_consumed": 0,
          "created_at": "2024-12-24T11:55:36Z",
          "updated_at": "2024-12-24T11:55:36Z"
        }
      ]
    }
  }
  ```

### 3.3 获取对话详情
- **接口**：`GET /conversations/detail/:id`
- **描述**：获取对话详情及其消息
- **请求头**：
  ```
  Authorization: Bearer <token>
  ```
- **响应**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "id": 1,
      "user_id": 1,
      "model_id": 1,
      "title": "测试对话",
      "points_consumed": 0,
      "created_at": "2024-12-24T11:55:36Z",
      "updated_at": "2024-12-24T11:55:36Z",
      "messages": [
        {
          "id": 1,
          "conversation_id": 1,
          "role": "user",
          "content": "你好",
          "created_at": "2024-12-24T11:56:00Z"
        }
      ]
    }
  }
  ```

### 3.4 添加消息
- **接口**：`POST /conversations/messages`
- **描述**：向对话添加新消息
- **请求头**：
  ```
  Content-Type: application/json
  Authorization: Bearer <token>
  ```
- **请求体**：
  ```json
  {
    "conversation_id": 1,  // 对话ID
    "role": "user",       // 角色：system/user/assistant
    "content": "你好"      // 消息内容
  }
  ```
- **响应**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "id": 1,
      "conversation_id": 1,
      "role": "user",
      "content": "你好",
      "created_at": "2024-12-24T11:56:00Z"
    }
  }
  ```

### 3.5 获取消息列表
- **接口**：`GET /conversations/messages/:conversation_id`
- **描述**：获取对话的消息列表
- **请求头**：
  ```
  Authorization: Bearer <token>
  ```
- **响应**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": [
      {
        "id": 1,
        "conversation_id": 1,
        "role": "user",
        "content": "你好",
        "created_at": "2024-12-24T11:56:00Z"
      }
    ]
  }
  ```

## 4. 错误码说明

| 错���码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 未授权 |
| 1003 | 禁止访问 |
| 1004 | 资源不存在 |
| 1005 | 系统错误 |

## 5. 测试用例

### 5.1 创建对话测试
```bash
curl -X POST http://localhost:8081/api/v1/conversations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJleHAiOjE3MDM1MjY0MDB9.K7DnXnJTZQxZ9ygPxr85c5rX9z7Y8sVJOvHJ5jQVTFE" \
  -d '{
    "model_id": 1,
    "title": "测试对话"
  }'
```

### 5.2 获取对话列表测试
```bash
curl -X GET "http://localhost:8081/api/v1/conversations?page=1&size=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJleHAiOjE3MDM1MjY0MDB9.K7DnXnJTZQxZ9ygPxr85c5rX9z7Y8sVJOvHJ5jQVTFE"
```

### 5.3 添加消息测试
```bash
curl -X POST http://localhost:8081/api/v1/conversations/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJleHAiOjE3MDM1MjY0MDB9.K7DnXnJTZQxZ9ygPxr85c5rX9z7Y8sVJOvHJ5jQVTFE" \
  -d '{
    "conversation_id": 1,
    "role": "user",
    "content": "你好"
  }'
```

### 5.4 获取消息列表测试
```bash
curl -X GET http://localhost:8081/api/v1/conversations/messages/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJleHAiOjE3MDM1MjY0MDB9.K7DnXnJTZQxZ9ygPxr85c5rX9z7Y8sVJOvHJ5jQVTFE"
```

## 6. 开发说明

### 6.1 目录结构
```
cybermind/chat-service/
├── cmd/
│   └── server/
│       └── main.go           # 程序入口
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   └── chat_handler.go  # API处理器
│   │   └── router/
│   │       └── router.go     # 路由配置
│   ├── model/
│   │   └── model.go         # 数据模型
│   └── service/
│       └── chat_service.go  # 业务逻辑
├── pkg/
│   ├── database/
│   │   └── database.go     # 数据库配置
│   └── middleware/
│       └── auth.go         # 认证中间件
└── go.mod                  # 依赖管理
```

### 6.2 认证说明
- 使用JWT进行身份认证
- Token格式：`Bearer <token>`
- Token中包含用户ID和过期时间
- Token密钥需要与auth-service保持一致

### 6.3 数据库说明
- 使用PostgreSQL存储数据
- 使用GORM进行数据库操作
- 自动创建数据表和索引
- 支持软删除

### 6.4 性能优化
- 使用数据库索引提升查询性能
- 分页查询避免大量数据返回
- 预加载关联数据减少查询次数
- 使用连接池管理数据库连接

## 7. 部署说明

### 7.1 环境要求
- Go 1.22或以上
- PostgreSQL 14.8.0或以上
- 操作系统：Linux/Windows

### 7.2 部署步骤
1. 克隆代码
2. 安装依赖：`go mod tidy`
3. 配置数据库连接
4. 编译：`go build cmd/server/main.go`
5. 运行：`./main`

### 7.3 配置说明
- 服务端口：8081
- 数据库连接：使用环境变量或配置文件
- JWT密钥：需要在环境变量中配置 