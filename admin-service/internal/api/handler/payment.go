package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"cybermind/admin-service/internal/model"
	"cybermind/admin-service/pkg/database"

	"github.com/gin-gonic/gin"
)

const (
	paymentDetailCacheKey = "payment:detail:%d"
	paymentCacheTTL      = 30 * time.Minute
)

// GetPaymentList 获取支付记录列表
func GetPaymentList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	orderNo := c.Query("order_no")
	paymentNo := c.Query("payment_no")
	status := c.Query("status")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	var payments []struct {
		ID            int64     `json:"id"`
		OrderNo       string    `json:"order_no"`
		PaymentNo     string    `json:"payment_no"`
		PaymentMethod string    `json:"payment_method"`
		Amount        float64   `json:"amount"`
		Status        int       `json:"status"`
		PaymentTime   time.Time `json:"payment_time"`
		CreatedAt     time.Time `json:"created_at"`
	}

	var total int64
	query := database.DB.Table("payments").
		Select("payments.*, orders.order_no").
		Joins("left join orders on orders.id = payments.order_id")

	if orderNo != "" {
		query = query.Where("orders.order_no = ?", orderNo)
	}
	if paymentNo != "" {
		query = query.Where("payments.payment_no = ?", paymentNo)
	}
	if status != "" {
		query = query.Where("payments.status = ?", status)
	}
	if startTime != "" {
		query = query.Where("payments.created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("payments.created_at <= ?", endTime)
	}

	query.Count(&total)
	if err := query.Offset((page - 1) * size).Limit(size).Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "获取成功",
		Data: model.PageResponse{
			Total: total,
			List:  payments,
		},
	})
}

// PaymentDetailResponse 支付详情响应
type PaymentDetailResponse struct {
	ID            int64     `json:"id" gorm:"column:id"`
	OrderID       int64     `json:"order_id" gorm:"column:order_id"`
	OrderNo       string    `json:"order_no" gorm:"column:order_no"`
	PaymentNo     string    `json:"payment_no" gorm:"column:payment_no"`
	PaymentMethod string    `json:"payment_method" gorm:"column:payment_method"`
	Amount        float64   `json:"amount" gorm:"column:amount"`
	Status        int       `json:"status" gorm:"column:status"`
	PaymentTime   time.Time `json:"payment_time" gorm:"column:payment_time"`
	RefundTime    time.Time `json:"refund_time" gorm:"column:refund_time"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at"`
	Callbacks     []model.PaymentCallback `json:"callbacks" gorm:"-"`
	Refunds       []model.PaymentRefund   `json:"refunds" gorm:"-"`
}

// GetPaymentDetail 获取支付记录详情
func GetPaymentDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	var payment PaymentDetailResponse

	// 尝试从缓存获取
	cacheKey := fmt.Sprintf(paymentDetailCacheKey, id)
	if cacheData, err := database.GetKey(c, cacheKey); err == nil {
		if err := json.Unmarshal([]byte(cacheData), &payment); err == nil {
			c.JSON(http.StatusOK, model.Response{
				Code:    model.Success,
				Message: "获取成功",
				Data:    payment,
			})
			return
		}
	}

	// 使用预加载获取所有相关数据
	if err := database.DB.Table("payments").
		Select("payments.*, orders.order_no").
		Joins("left join orders on orders.id = payments.order_id").
		Where("payments.id = ?", id).
		Take(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{
			Code:    model.NotFound,
			Message: "支付记录不存在",
		})
		return
	}

	// 获取回调记录和退款记录
	if err := database.DB.Model(&model.PaymentCallback{}).Where("payment_id = ?", id).Find(&payment.Callbacks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	if err := database.DB.Model(&model.PaymentRefund{}).Where("payment_id = ?", id).Find(&payment.Refunds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	// 尝试缓存数据
	if cacheData, err := json.Marshal(payment); err == nil {
		if err := database.SetKey(c, cacheKey, string(cacheData), paymentCacheTTL); err != nil {
			log.Printf("Failed to set cache: %v", err)
		}
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "获取成功",
		Data:    payment,
	})
}

// CreateRefundRequest 创建退款请求
type CreateRefundRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Reason string  `json:"reason" binding:"required"`
}

// CreateRefund 创建退款
func CreateRefund(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	var req CreateRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	// 开始事务
	tx := database.DB.Begin()

	// 检查支付记录
	var payment model.Payment
	if err := tx.First(&payment, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, model.Response{
			Code:    model.NotFound,
			Message: "支付记录不存在",
		})
		return
	}

	// 检查支付状态
	if payment.Status != 1 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "支付状态不正确",
		})
		return
	}

	// 检查退款金额
	if req.Amount > payment.Amount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "退款金额不能大于支付金额",
		})
		return
	}

	// 创建退款记录
	refund := &model.PaymentRefund{
		PaymentID: id,
		RefundNo:  "RF" + time.Now().Format("20060102150405") + strconv.FormatInt(id, 10),
		Amount:    req.Amount,
		Reason:    req.Reason,
		Status:    0,
	}
	if err := tx.Create(refund).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	// 更新支付状态
	if err := tx.Model(&payment).Updates(map[string]interface{}{
		"status":      3,
		"refund_time": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	// 记录操作日志
	adminID, exists := c.Get("admin_id")
	if !exists {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	operation := &model.AdminOperation{
		AdminID:     adminID.(int64),
		Module:      "payment",
		Action:      "create_refund",
		Description: req.Reason,
		IP:          c.ClientIP(),
	}
	if err := tx.Create(operation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	// 尝试删除缓存
	cacheKey := fmt.Sprintf(paymentDetailCacheKey, id)
	if err := database.DelKey(context.Background(), cacheKey); err != nil {
		log.Printf("Failed to delete cache: %v", err)
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "创建成功",
		Data:    refund,
	})
} 