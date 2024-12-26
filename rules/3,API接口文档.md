# CyberMind Π API接口文档

## 1. 接口规范

### 1.1 请求规范
- 基础URL：`/api/v1`
- 请求方法：GET、POST、PUT、DELETE
- 请求头：
  ```
  Content-Type: application/json
  Authorization: Bearer <token>
  ```

### 1.2 响应规范
```json
{
    "code": 0,           // 状态码：0成功，非0失败
    "message": "success", // 状态信息
    "data": {}           // 响应数据
}
```

### 1.3 错误码说明
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
| 3001 | 积分不足 |
| 3002 | 套餐已过期 |
| 3003 | 支付失败 |

## 2. 用户相关接口

### 2.1 用户注册
```http
POST /api/v1/auth/register

Request:
{
    "username": "string",     // 用户名
    "email": "string",        // 邮箱
    "password": "string",     // 密码
    "phone": "string"         // 手机号（可选）
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "user_id": "int64",
        "username": "string",
        "email": "string"
    }
}
```

### 2.2 用户登录
```http
POST /api/v1/auth/login

Request:
{
    "email": "string",        // 邮箱
    "password": "string"      // 密码
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "token": "string",
        "user": {
            "id": "int64",
            "username": "string",
            "email": "string",
            "points": "int"
        }
    }
}
```

### 2.3 获取用户信息
```http
GET /api/v1/users/info

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "id": "int64",
        "username": "string",
        "email": "string",
        "phone": "string",
        "points": "int",
        "avatar_url": "string"
    }
}
```

### 2.4 更新用户信息
```http
PUT /api/v1/users/info

Request:
{
    "username": "string",     // 用户名（可选）
    "phone": "string",        // 手机号（可选）
    "avatar_url": "string"    // 头像URL（可选）
}

Response:
{
    "code": 0,
    "message": "success",
    "data": null
}
```

## 3. 对话相关接口

### 3.1 创建对话
```http
POST /api/v1/conversations

Request:
{
    "model_id": "int64",
    "messages": [
        {
            "role": "string",     // system/user/assistant
            "content": "string"
        }
    ]
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "conversation_id": "int64",
        "message": {
            "role": "string",
            "content": "string"
        },
        "points_consumed": "int"
    }
}
```

### 3.2 获取对话历史
```http
GET /api/v1/conversations?page=1&size=20

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "total": "int",
        "list": [
            {
                "id": "int64",
                "title": "string",
                "model_id": "int64",
                "points_consumed": "int",
                "created_at": "string"
            }
        ]
    }
}
```

### 3.3 获取对话详情
```http
GET /api/v1/conversations/{id}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "id": "int64",
        "title": "string",
        "model_id": "int64",
        "messages": [
            {
                "role": "string",
                "content": "string",
                "created_at": "string"
            }
        ],
        "points_consumed": "int",
        "created_at": "string"
    }
}
```

## 4. 套餐相关接口

### 4.1 获取套餐列表
```http
GET /api/v1/packages

Response:
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "id": "int64",
            "name": "string",
            "type": "int",
            "price": "float",
            "total_points": "int",
            "daily_points": "int",
            "duration_days": "int"
        }
    ]
}
```

### 4.2 创建订单
```http
POST /api/v1/orders

Request:
{
    "package_id": "int64"
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "order_no": "string",
        "amount": "float",
        "pay_url": "string"
    }
}
```

### 4.3 查询订单状态
```http
GET /api/v1/orders/{order_no}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "order_no": "string",
        "status": "int",
        "amount": "float",
        "payment_time": "string"
    }
}
```

## 5. 管理员接口

### 5.1 用户管理

#### 5.1.1 获取用户列表
```http
GET /api/v1/admin/users?page=1&size=20

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "total": "int",
        "list": [
            {
                "id": "int64",
                "username": "string",
                "email": "string",
                "status": "int",
                "points": "int",
                "created_at": "string"
            }
        ]
    }
}
```

#### 5.1.2 更新用户状态
```http
PUT /api/v1/admin/users/{id}/status

Request:
{
    "status": "int",      // 0禁用，1正常
    "reason": "string"    // 操作原因
}

Response:
{
    "code": 0,
    "message": "success",
    "data": null
}
```

### 5.2 模型管理

#### 5.2.1 获取模型列表
```http
GET /api/v1/admin/models

Response:
{
    "code": 0,
    "message": "success",
    "data": [
        {
            "id": "int64",
            "name": "string",
            "provider": "string",
            "api_type": "string",
            "points_per_request": "int",
            "status": "int",
            "config": {}
        }
    ]
}
```

#### 5.2.2 更新模型配置
```http
PUT /api/v1/admin/models/{id}

Request:
{
    "name": "string",
    "points_per_request": "int",
    "status": "int",
    "config": {
        "temperature": "float",
        "max_tokens": "int",
        "top_p": "float",
        "frequency_penalty": "float",
        "presence_penalty": "float"
    }
}

Response:
{
    "code": 0,
    "message": "success",
    "data": null
}
```

### 5.3 统计分析

#### 5.3.1 获取系统概览
```http
GET /api/v1/admin/dashboard/overview

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "total_users": "int",
        "active_users": "int",
        "total_conversations": "int",
        "total_orders": "int",
        "total_amount": "float"
    }
}
```

#### 5.3.2 获取使用统计
```http
GET /api/v1/admin/dashboard/stats?type=daily&start_date=string&end_date=string

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "dates": ["string"],
        "conversations": ["int"],
        "new_users": ["int"],
        "orders": ["int"],
        "amounts": ["float"]
    }
}
``` 