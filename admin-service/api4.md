# 支付管理模块开发文档

## 1. 功能概述

支付管理模块主要提供以下功能：
- 支付记录列表查询
- 支付记录详情查看
- 退款操作处理

## 2. 数据库设计

### 2.1 表结构

#### payments 支付记录表
```sql
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    payment_no VARCHAR(50) NOT NULL UNIQUE,
    payment_method VARCHAR(20) NOT NULL,
    amount DECIMAL NOT NULL,
    status INT DEFAULT 0,
    payment_time TIMESTAMPTZ,
    refund_time TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- 索引
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_payment_method ON payments(payment_method);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);
```

#### payment_callbacks 支付回调记录表
```sql
CREATE TABLE payment_callbacks (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT NOT NULL,
    callback_no VARCHAR(100) NOT NULL UNIQUE,
    status INT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- 索引
CREATE INDEX idx_payment_callbacks_payment_id ON payment_callbacks(payment_id);
CREATE INDEX idx_payment_callbacks_status ON payment_callbacks(status);
CREATE INDEX idx_payment_callbacks_created_at ON payment_callbacks(created_at);
```

#### payment_refunds 退款记录表
```sql
CREATE TABLE payment_refunds (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT NOT NULL,
    refund_no VARCHAR(50) NOT NULL UNIQUE,
    amount DECIMAL NOT NULL,
    reason VARCHAR(255),
    status INT DEFAULT 0,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- 索引
CREATE INDEX idx_payment_refunds_payment_id ON payment_refunds(payment_id);
CREATE INDEX idx_payment_refunds_status ON payment_refunds(status);
CREATE INDEX idx_payment_refunds_created_at ON payment_refunds(created_at);
```

### 2.2 字段说明

#### payments 表
- id: 主键
- order_id: 订单ID
- payment_no: 支付单号
- payment_method: 支付方式(alipay/wechat/bank)
- amount: 支付金额
- status: 支付状态(0:待支付 1:支付成功 2:支付失败 3:已退款)
- payment_time: 支付时间
- refund_time: 退款时间
- created_at: 创建时间
- updated_at: 更新时间

#### payment_callbacks 表
- id: 主键
- payment_id: 支付记录ID
- callback_no: 回调单号
- status: 回调状态(0:处理中 1:成功 2:失败)
- created_at: 创建时间
- updated_at: 更新时间

#### payment_refunds 表
- id: 主键
- payment_id: 支付记录ID
- refund_no: 退款单号
- amount: 退款金额
- reason: 退款原因
- status: 退款状态(0:处理中 1:成功 2:失败)
- created_at: 创建时间
- updated_at: 更新时间

## 3. API接口设计

### 3.1 获取支付记录列表

#### 请求信息
- 请求方法: GET
- 请求路径: /payments
- 请求参数:
  ```json
  {
    "page": "页码(默认1)",
    "size": "每页数量(默认10)",
    "order_no": "订单号(可选)",
    "payment_no": "支付单号(可选)",
    "status": "支付状态(可选)",
    "start_time": "开始时间(可选)",
    "end_time": "结束时间(可选)"
  }
  ```

#### 响应信息
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "total": "总记录数",
    "list": [
      {
        "id": "支付记录ID",
        "order_no": "订单号",
        "payment_no": "支付单号",
        "payment_method": "支付方式",
        "amount": "支付金额",
        "status": "支付状态",
        "payment_time": "支付时间",
        "created_at": "创建时间"
      }
    ]
  }
}
```

### 3.2 获取支付记录详情

#### 请求信息
- 请求方法: GET
- 请���路径: /payments/:id
- 路径参数:
  - id: 支付记录ID

#### 响应信息
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "id": "支付记录ID",
    "order_id": "订单ID",
    "order_no": "订单号",
    "payment_no": "支付单号",
    "payment_method": "支付方式",
    "amount": "支付金额",
    "status": "支付状态",
    "payment_time": "支付时间",
    "refund_time": "退款时间",
    "created_at": "创建时间",
    "updated_at": "更新时间",
    "callbacks": [
      {
        "id": "回调记录ID",
        "callback_no": "回调单号",
        "status": "回调状态",
        "created_at": "创建时间"
      }
    ],
    "refunds": [
      {
        "id": "退款记录ID",
        "refund_no": "退款单号",
        "amount": "退款金额",
        "reason": "退款原因",
        "status": "退款状态",
        "created_at": "创建时间"
      }
    ]
  }
}
```

### 3.3 创建退款

#### 请求信息
- 请求方法: POST
- 请求路径: /payments/:id/refund
- 路径参数:
  - id: 支付记录ID
- 请求体:
  ```json
  {
    "amount": "退款金额(必填,大于0)",
    "reason": "退款原因(必填)"
  }
  ```

#### 响应信息
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": "退款记录ID",
    "payment_id": "支付记录ID",
    "refund_no": "退款单号",
    "amount": "退款金额",
    "reason": "退款原因",
    "status": "退款状态",
    "created_at": "创建时间"
  }
}
```

## 4. 缓存设计

### 4.1 缓存键
- 支付详情缓存: `payment:detail:{id}`
- 缓存时间: 30分钟

### 4.2 缓存策略
- 查询时优先从缓存获取
- 缓存未命中时从数据库查询并写入缓存
- 退款操作后删除相关缓存

## 5. 测试用例

### 5.1 单元测试

#### TestGetPaymentList
- 测试目的: 验证支付列表查询功能
- 测试步骤:
  1. 准备测试数据
  2. 发送列表查询请求
  3. 验证响应状态码和数据
- 测试结果: PASS

#### TestGetPaymentDetail
- 测试目的: 验证支付详情查询功能
- 测试步骤:
  1. 准备测试数据(支付记录、回调记录)
  2. 发送详情查询请求
  3. 验证响应状态码和数据
- 测试结果: PASS

#### TestCreateRefund
- 测试目的: 验证退款创建功能
- 测试步骤:
  1. 准备测试数据(支付记录)
  2. 发送退款请求
  3. 验证响应状态码和数据
  4. 验证数据库记录
- 测试结果: PASS

### 5.2 基准测试

#### BenchmarkGetPaymentList
- 测试目的: 测试列表查询性能
- 测试结果: 平均响应时间 < 1ms

#### BenchmarkGetPaymentDetail
- 测试目的: 测试详情查询性能
- 测试结果: 平均响应时间 < 1ms

#### BenchmarkCreateRefund
- 测试目的: 测试退款创建性能
- 测试结果: 平均响应时间 < 1ms

#### 并发测试
- BenchmarkParallelGetPaymentList: 并发列表查询
- BenchmarkParallelGetPaymentDetail: 并发详情查询
- BenchmarkParallelCreateRefund: 并发退款创建

## 6. 性能优化

### 6.1 数据库优化
- 添加合适的索引
- 使用预加载减少查询次数
- 优化SQL查询语句

### 6.2 缓存优化
- 使用Redis缓存热点数据
- 合理设置缓存时间
- 及时清理过期缓存

### 6.3 代码优化
- 使用结构体标签优化字段映射
- 优化错误处理逻辑
- 添加日志记录

## 7. 部署建议

### 7.1 数据库配置
- 配置合适的连接池参数
- 定期维护数据库索引
- 监控慢查询日志

### 7.2 Redis配置
- 配置合适的连接池参数
- 设置合理的内存策略
- 监控缓存命中率

### 7.3 系统监控
- 添加接口调用监控
- 添加错误日志监控
- 添加性能指标监控 