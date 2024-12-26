# CyberMind认证服务API文档

## 1. 服务概述

### 1.1 基本信息
- 服务名称：认证服务（Auth Service）
- 基础URL：`http://localhost:8081/api/v1`
- 开发语言：Go
- 主要依赖：
  - gin-gonic/gin: Web框架
  - gorm.io/gorm: ORM框架
  - golang-jwt/jwt: JWT认证
  - golang.org/x/crypto: 密码加密

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

### 2.1 用户模型（User）
```go
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
```

## 3. API接口

### 3.1 用户注册
- **接口**：`POST /auth/register`
- **描述**：注册新用户
- **请求体**：
  ```json
  {
    "username": "string",     // 用户名（3-50字符）
    "email": "string",        // 邮箱（有效的邮箱格式）
    "password": "string",     // 密码（最少8位）
    "phone": "string"         // 手机号（11位）
  }
  ```
- **响应**：
  ```json
  {
    "code": 0,               // 状态码：0成功，非0失败
    "message": "success",    // 状态信息
    "data": {
      "user_id": 1,         // 用户ID
      "username": "string", // 用户名
      "email": "string"     // 邮箱
    }
  }
  ```
- **错误码**：
  - 1001：参数错误或用户已存在
  - 1005：系统错误

### 3.2 用户登录
- **接口**：`POST /auth/login`
- **描述**：用户登录并获取JWT令牌
- **请求体**：
  ```json
  {
    "email": "string",    // 邮箱
    "password": "string"  // 密码
  }
  ```
- **响应**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "token": "string",           // JWT令牌
      "user": {
        "id": 1,                   // 用户ID
        "username": "string",      // 用户名
        "email": "string",         // 邮箱
        "points": 0                // 积分
      }
    }
  }
  ```
- **错误码**：
  - 1001：参数错误
  - 1002：用户不存在或密码错误
  - 2003：账号已禁用

### 3.3 获取用户信息
- **接口**：`GET /users/info`
- **描述**：获取当前登录用户的详细信息
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
      "id": 1,                // 用户ID
      "username": "string",   // 用户名
      "email": "string",      // 邮箱
      "phone": "string",      // 手机号
      "points": 0,            // 积分
      "avatar_url": "string"  // 头像URL
    }
  }
  ```
- **错误码**：
  - 1002：未授权
  - 1004：用户不存在

## 4. 安全特性

### 4.1 密码安全
- 使用bcrypt算法加密存储密码
- 密码最小长度要求：8位
- 密码从不明文传输和存储

### 4.2 JWT认证
- Token默认有效期：24小时
- Token包含信息：
  - 用户ID
  - 用户名
  - 用户角色
  - 过期时间

### 4.3 权限控制
- 用户角色：
  - 0：普通用户
  - 1：管理员（特殊邮箱：2986309418@qq.com）
- 用户状态：
  - 1：正常
  - 0：禁用

## 5. 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 未授权 |
| 1003 | 禁止访问 |
| 1004 | 资源不存在 |
| 1005 | 系统错误 |
| 2001 | 用户不存在 |
| 2002 | 密码错误 |
| 2003 | 账号已禁用 |

## 6. 测试用例

### 6.1 注册测试
```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test123456",
    "phone": "13800138000"
  }'
```

### 6.2 登录测试
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123456"
  }'
```

### 6.3 获取用户信息测试
```bash
curl -X GET http://localhost:8081/api/v1/users/info \
  -H "Authorization: Bearer <token>"
```

## 7. 部署信息

### 7.1 环境要求
- Go 1.22或以上
- PostgreSQL 14.8.0或以上
- 操作系统：Linux/Windows

### 7.2 配置说明
- 服务端口：8081
- 数据库连接：使用环境变量或配置文件
- JWT密钥：需要在环境变量中配置

### 7.3 启动命令
```bash
cd /home/devbox/project/cybermind/auth-service
go run cmd/server/main.go
``` 