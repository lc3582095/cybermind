package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cybermind/admin-service/internal/api/handler"
	"cybermind/admin-service/internal/model"
	"cybermind/admin-service/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(postgres.Open("host=dbconn.sealosbja.site port=37550 user=postgres password=wkzhx7jn dbname=postgres sslmode=disable"), &gorm.Config{})
	assert.NoError(t, err)

	// 清理测试数据
	database.DB.Exec("DROP TABLE IF EXISTS payments, payment_callbacks, payment_refunds, orders, admin_operations")
	
	// 创建 orders 表
	database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id BIGSERIAL PRIMARY KEY,
			order_no VARCHAR(50) NOT NULL UNIQUE,
			user_id BIGINT NOT NULL,
			amount DECIMAL NOT NULL,
			status INT DEFAULT 0,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		)
	`)

	// 创建 admin_operations 表
	database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS admin_operations (
			id BIGSERIAL PRIMARY KEY,
			admin_id BIGINT NOT NULL,
			module VARCHAR(50) NOT NULL,
			action VARCHAR(50) NOT NULL,
			description TEXT,
			ip VARCHAR(50),
			created_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)

	err = database.DB.AutoMigrate(&model.Payment{}, &model.PaymentCallback{}, &model.PaymentRefund{})
	assert.NoError(t, err)

	// 创建测试订单
	database.DB.Exec(`
		INSERT INTO orders (order_no, user_id, amount, status, created_at)
		VALUES ('TEST_ORDER_001', 1, 100, 1, NOW())
	`)
}

func TestGetPaymentList(t *testing.T) {
	setupTestDB(t)

	// 准备测试数据
	payment := &model.Payment{
		OrderID:       1,
		PaymentNo:     "TEST001",
		PaymentMethod: "alipay",
		Amount:        100,
		Status:        1,
		PaymentTime:   time.Now(),
		CreatedAt:     time.Now(),
	}
	err := database.DB.Create(payment).Error
	assert.NoError(t, err)

	// 设置测试路由
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/payments", handler.GetPaymentList)

	// 发送测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/payments", nil)
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    model.PageResponse `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, model.Success, resp.Code)
	assert.Equal(t, int64(1), resp.Data.Total)
}

func TestGetPaymentDetail(t *testing.T) {
	setupTestDB(t)

	// 准备测试数据
	payment := &model.Payment{
		OrderID:       1,
		PaymentNo:     "TEST001",
		PaymentMethod: "alipay",
		Amount:        100,
		Status:        1,
		PaymentTime:   time.Now(),
		CreatedAt:     time.Now(),
	}
	err := database.DB.Create(payment).Error
	assert.NoError(t, err)

	callback := &model.PaymentCallback{
		PaymentID:  payment.ID,
		CallbackNo: "CB001",
		Status:     1,
		CreatedAt:  time.Now(),
	}
	err = database.DB.Create(callback).Error
	assert.NoError(t, err)

	// 设置测试路由
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/payments/:id", handler.GetPaymentDetail)

	// 发送测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/payments/1", nil)
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, model.Success, resp.Code)
}

func TestCreateRefund(t *testing.T) {
	setupTestDB(t)

	// 准备测试数据
	payment := &model.Payment{
		OrderID:       1,
		PaymentNo:     "TEST001",
		PaymentMethod: "alipay",
		Amount:        100,
		Status:        1,
		PaymentTime:   time.Now(),
		CreatedAt:     time.Now(),
	}
	err := database.DB.Create(payment).Error
	assert.NoError(t, err)

	// 设置测试路由
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("admin_id", int64(1))
		c.Next()
	})
	r.POST("/payments/:id/refund", handler.CreateRefund)

	// 准备请求数据
	reqBody := handler.CreateRefundRequest{
		Amount: 50,
		Reason: "测试退款",
	}
	reqData, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	// 发送测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/payments/1/refund", bytes.NewBuffer(reqData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, model.Success, resp.Code)

	// 验证数据库记录
	var refund model.PaymentRefund
	err = database.DB.Where("payment_id = ?", payment.ID).First(&refund).Error
	assert.NoError(t, err)
	assert.Equal(t, float64(50), refund.Amount)
	assert.Equal(t, "测试退款", refund.Reason)
} 