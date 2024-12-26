package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cybermind/admin-service/internal/model"
	"cybermind/admin-service/pkg/database"

	"github.com/gin-gonic/gin"
)

// GetOrderList 获取订单列表
func GetOrderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	orderNo := c.Query("order_no")
	status := c.Query("status")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	var orders []struct {
		ID        int64     `json:"id"`
		OrderNo   string    `json:"order_no"`
		UserID    int64     `json:"user_id"`
		Amount    float64   `json:"amount"`
		Status    int       `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	var total int64
	query := database.DB.Table("orders")

	if orderNo != "" {
		query = query.Where("order_no = ?", orderNo)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	query.Count(&total)
	if err := query.Offset((page - 1) * size).Limit(size).Find(&orders).Error; err != nil {
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
			List:  orders,
		},
	})
}

// GetOrderDetail 获取订单详情
func GetOrderDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	var order struct {
		ID        int64     `json:"id"`
		OrderNo   string    `json:"order_no"`
		UserID    int64     `json:"user_id"`
		Amount    float64   `json:"amount"`
		Status    int       `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		Items     []struct {
			ID       int64   `json:"id"`
			OrderID  int64   `json:"order_id"`
			PackageID int64  `json:"package_id"`
			Points   int     `json:"points"`
			Amount   float64 `json:"amount"`
		} `json:"items"`
		Payments []struct {
			ID            int64     `json:"id"`
			PaymentNo     string    `json:"payment_no"`
			PaymentMethod string    `json:"payment_method"`
			Amount        float64   `json:"amount"`
			Status        int       `json:"status"`
			PaymentTime   time.Time `json:"payment_time"`
		} `json:"payments"`
	}

	// 获取订单信息
	if err := database.DB.Table("orders").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{
			Code:    model.NotFound,
			Message: "订单不存在",
		})
		return
	}

	// 获取订单项
	database.DB.Table("order_items").
		Where("order_id = ?", id).
		Find(&order.Items)

	// 获取支付记录
	database.DB.Table("payments").
		Where("order_id = ?", id).
		Find(&order.Payments)

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "获取成功",
		Data:    order,
	})
}

// UpdateOrderStatusRequest 更新订单状态请求
type UpdateOrderStatusRequest struct {
	Status int    `json:"status" binding:"required,oneof=0 1 2 3"`
	Reason string `json:"reason" binding:"required"`
}

// UpdateOrderStatus 更新订单状态
func UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    model.ParamError,
			Message: "参数错误",
		})
		return
	}

	// 开始事务
	tx := database.DB.Begin()

	// 检查订单是否存在
	var order struct {
		ID     int64 `json:"id"`
		Status int   `json:"status"`
	}
	if err := tx.Table("orders").First(&order, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, model.Response{
			Code:    model.NotFound,
			Message: "订单不存在",
		})
		return
	}

	// 更新订单状态
	if err := tx.Table("orders").Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    model.SystemError,
			Message: "系统错误",
		})
		return
	}

	// 记录操作日志
	adminID, _ := c.Get("admin_id")
	operation := &model.AdminOperation{
		AdminID:     adminID.(int64),
		Module:      "order",
		Action:      fmt.Sprintf("update_status_%d", req.Status),
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

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "更新成功",
	})
}