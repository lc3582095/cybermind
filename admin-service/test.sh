#!/bin/bash

# 设置基础URL
BASE_URL="http://localhost:8081/api/v1"
TOKEN=""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 测试计数器
TOTAL=0
PASSED=0

# 测试函数
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_code=$5
    
    TOTAL=$((TOTAL + 1))
    
    echo -e "\n测试用例: ${name}"
    echo "请求: ${method} ${endpoint}"
    if [ ! -z "$data" ]; then
        echo "数据: ${data}"
    fi
    
    # 构建curl命令
    local curl_cmd="curl -s -X ${method} '${BASE_URL}${endpoint}'"
    if [ ! -z "$data" ]; then
        curl_cmd="${curl_cmd} -H 'Content-Type: application/json' -d '${data}'"
    fi
    if [ ! -z "$TOKEN" ]; then
        curl_cmd="${curl_cmd} -H 'Authorization: Bearer ${TOKEN}'"
    fi
    
    # 执行请求
    local response=$(eval ${curl_cmd})
    local actual_code=$(echo ${response} | jq -r '.code')
    
    # 检查响应
    if [ "${actual_code}" = "${expected_code}" ]; then
        echo -e "${GREEN}通过√${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}失败×${NC}"
        echo "期望状态码: ${expected_code}"
        echo "实际状态码: ${actual_code}"
        echo "响应: ${response}"
    fi
}

# 1. 认证测试
echo "=== 认证测试 ==="

# 1.1 登录测试 - 成功
test_api "管理员登录 - 成功" "POST" "/auth/login" \
    '{"email":"2986309418@qq.com","password":"admin123"}' "0"

# 保存token
TOKEN=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"2986309418@qq.com","password":"admin123"}' \
    | jq -r '.data.token')

# 1.2 登录测试 - 失败（密码错误）
test_api "管理员登录 - 密码错误" "POST" "/auth/login" \
    '{"email":"2986309418@qq.com","password":"wrong"}' "2002"

# 1.3 获取管理员信息
test_api "获取管理员信息" "GET" "/auth/info" "" "0"

# 2. 管理员管理测试
echo -e "\n=== 管理员管理测试 ==="

# 2.1 创建管理员
test_api "创建管理员" "POST" "/admin/create" \
    '{"username":"test_admin","email":"test@example.com","password":"123456","role":1}' "0"

# 2.2 获取管理员列表
test_api "获取管理员列表" "GET" "/admin/list" "" "0"

# 2.3 更新管理员
test_api "更新管理员" "PUT" "/admin/2" \
    '{"username":"test_admin_updated","role":1,"status":1}' "0"

# 3. 用户管理测试
echo -e "\n=== 用户管理测试 ==="

# 3.1 获取用户列表
test_api "获取用户列表" "GET" "/admin/users" "" "0"

# 3.2 获取用户详情
test_api "获取用户详情" "GET" "/admin/users/1/detail" "" "0"

# 3.3 更新用户状态
test_api "更新用户状态" "PUT" "/admin/users/1/status" \
    '{"status":1,"reason":"测试"}' "0"

# 4. 模型管理测试
echo -e "\n=== 模型管理测试 ==="

# 4.1 创建模型
test_api "创建模型" "POST" "/admin/models" \
    '{"name":"GPT-4","provider":"OpenAI","api_type":"chat","base_url":"https://api.openai.com","api_key":"sk-xxx","model_name":"gpt-4","points_per_request":10}' "0"

# 4.2 获取模型列表
test_api "获取模型列表" "GET" "/admin/models" "" "0"

# 4.3 更新模型
test_api "更新模型" "PUT" "/admin/models/1" \
    '{"name":"GPT-4","base_url":"https://api.openai.com","api_key":"sk-xxx","points_per_request":12,"status":1}' "0"

# 5. 订单管理测试
echo -e "\n=== 订单管理测试 ==="

# 5.1 获取订单列表
test_api "获取订单列表" "GET" "/admin/orders" "" "0"

# 5.2 获取订单详情
test_api "获取订单详情" "GET" "/admin/orders/1" "" "0"

# 5.3 更新订单状态
test_api "更新订单状态" "PUT" "/admin/orders/1/status" \
    '{"status":1,"reason":"测试"}' "0"

# 6. 支付管理测试
echo -e "\n=== 支付管理测试 ==="

# 6.1 获取支付列表
test_api "获取支付列表" "GET" "/admin/payments" "" "0"

# 6.2 获取支付详情
test_api "获取支付详情" "GET" "/admin/payments/1" "" "0"

# 6.3 创建退款
test_api "创建退款" "POST" "/admin/payments/1/refund" \
    '{"amount":10,"reason":"测试退款"}' "0"

# 7. 统计分析测试
echo -e "\n=== 统计分析测试 ==="

# 7.1 获取系统概览
test_api "获取系统概览" "GET" "/admin/stats/overview" "" "0"

# 7.2 获取每日统计
test_api "获取每日统计" "GET" "/admin/stats/daily" "" "0"

# 7.3 获取用户统计
test_api "获取用户统计" "GET" "/admin/stats/user" "" "0"

# 7.4 获取订单统计
test_api "获取订单统计" "GET" "/admin/stats/order" "" "0"

# 输出测试结果
echo -e "\n=== 测试结果 ==="
echo "总用例数: ${TOTAL}"
echo "通过数量: ${PASSED}"
echo "失败数量: $((TOTAL - PASSED))"
echo "通过率: $(( (PASSED * 100) / TOTAL ))%" 