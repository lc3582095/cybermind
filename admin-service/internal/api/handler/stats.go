package handler

import (
	"net/http"
	"time"

	"cybermind/admin-service/internal/model"
	"cybermind/admin-service/pkg/database"

	"github.com/gin-gonic/gin"
)

// GetStatsOverview 获取系统概览
func GetStatsOverview(c *gin.Context) {
	var stats struct {
		TotalUsers     int64   `json:"total_users"`
		ActiveUsers    int64   `json:"active_users"`
		TotalOrders    int64   `json:"total_orders"`
		TotalAmount    float64 `json:"total_amount"`
		TotalRequests  int64   `json:"total_requests"`
		SuccessRate    float64 `json:"success_rate"`
		AverageLatency float64 `json:"average_latency"`
	}

	// 获取用户统计
	database.DB.Table("users").Count(&stats.TotalUsers)
	database.DB.Table("users").Where("last_active_at >= ?", time.Now().AddDate(0, 0, -7)).Count(&stats.ActiveUsers)

	// 获取订单统计
	database.DB.Table("orders").Count(&stats.TotalOrders)
	database.DB.Table("orders").Select("COALESCE(SUM(amount), 0)").Row().Scan(&stats.TotalAmount)

	// 获取请求统计
	database.DB.Table("model_requests").Count(&stats.TotalRequests)
	database.DB.Table("model_requests").
		Select("COALESCE(AVG(CASE WHEN status = 1 THEN 1 ELSE 0 END), 0)").
		Row().Scan(&stats.SuccessRate)
	database.DB.Table("model_requests").
		Select("COALESCE(AVG(latency), 0)").
		Where("status = 1").
		Row().Scan(&stats.AverageLatency)

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "获取成功",
		Data:    stats,
	})
}

// GetDailyStats 获取每日统计
func GetDailyStats(c *gin.Context) {
	var stats []struct {
		Date           string  `json:"date"`
		NewUsers       int64   `json:"new_users"`
		ActiveUsers    int64   `json:"active_users"`
		OrderCount     int64   `json:"order_count"`
		OrderAmount    float64 `json:"order_amount"`
		RequestCount   int64   `json:"request_count"`
		SuccessRate    float64 `json:"success_rate"`
		AverageLatency float64 `json:"average_latency"`
	}

	database.DB.Raw(`
		WITH dates AS (
			SELECT generate_series(
				date_trunc('day', now()) - interval '29 days',
				date_trunc('day', now()),
				interval '1 day'
			)::date AS date
		)
		SELECT 
			d.date::text,
			COUNT(DISTINCT CASE WHEN u.created_at::date = d.date THEN u.id END) as new_users,
			COUNT(DISTINCT CASE WHEN u.last_active_at::date = d.date THEN u.id END) as active_users,
			COUNT(DISTINCT CASE WHEN o.created_at::date = d.date THEN o.id END) as order_count,
			COALESCE(SUM(CASE WHEN o.created_at::date = d.date THEN o.amount ELSE 0 END), 0) as order_amount,
			COUNT(DISTINCT CASE WHEN r.created_at::date = d.date THEN r.id END) as request_count,
			COALESCE(AVG(CASE WHEN r.created_at::date = d.date THEN CASE WHEN r.status = 1 THEN 1 ELSE 0 END END), 0) as success_rate,
			COALESCE(AVG(CASE WHEN r.created_at::date = d.date AND r.status = 1 THEN r.latency END), 0) as average_latency
		FROM dates d
		LEFT JOIN users u ON TRUE
		LEFT JOIN orders o ON TRUE
		LEFT JOIN model_requests r ON TRUE
		GROUP BY d.date
		ORDER BY d.date DESC
	`).Scan(&stats)

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "获取成功",
		Data:    stats,
	})
}

// GetUserStats 获取用户统计
func GetUserStats(c *gin.Context) {
	var stats struct {
		UserGrowth []struct {
			Date      string `json:"date"`
			NewUsers  int64  `json:"new_users"`
			TotalUsers int64 `json:"total_users"`
		} `json:"user_growth"`
		UserActivity []struct {
			Date        string `json:"date"`
			ActiveUsers int64  `json:"active_users"`
			RequestCount int64 `json:"request_count"`
		} `json:"user_activity"`
		UserPoints []struct {
			Points    int   `json:"points"`
			UserCount int64 `json:"user_count"`
		} `json:"user_points"`
	}

	// 获取用户增长数据
	database.DB.Raw(`
		WITH dates AS (
			SELECT generate_series(
				date_trunc('day', now()) - interval '29 days',
				date_trunc('day', now()),
				interval '1 day'
			)::date AS date
		)
		SELECT 
			d.date::text,
			COUNT(DISTINCT CASE WHEN u.created_at::date = d.date THEN u.id END) as new_users,
			COUNT(DISTINCT CASE WHEN u.created_at::date <= d.date THEN u.id END) as total_users
		FROM dates d
		LEFT JOIN users u ON TRUE
		GROUP BY d.date
		ORDER BY d.date DESC
	`).Scan(&stats.UserGrowth)

	// 获取用户活跃数据
	database.DB.Raw(`
		WITH dates AS (
			SELECT generate_series(
				date_trunc('day', now()) - interval '29 days',
				date_trunc('day', now()),
				interval '1 day'
			)::date AS date
		)
		SELECT 
			d.date::text,
			COUNT(DISTINCT CASE WHEN u.last_active_at::date = d.date THEN u.id END) as active_users,
			COUNT(DISTINCT CASE WHEN r.created_at::date = d.date THEN r.id END) as request_count
		FROM dates d
		LEFT JOIN users u ON TRUE
		LEFT JOIN model_requests r ON r.user_id = u.id
		GROUP BY d.date
		ORDER BY d.date DESC
	`).Scan(&stats.UserActivity)

	// 获取用户积分分布
	database.DB.Raw(`
		SELECT 
			CASE 
				WHEN points = 0 THEN 0
				WHEN points <= 100 THEN 100
				WHEN points <= 500 THEN 500
				WHEN points <= 1000 THEN 1000
				WHEN points <= 5000 THEN 5000
				ELSE 10000
			END as points,
			COUNT(*) as user_count
		FROM users
		GROUP BY 1
		ORDER BY 1
	`).Scan(&stats.UserPoints)

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "获取成功",
		Data:    stats,
	})
}

// GetOrderStats 获取订单统计
func GetOrderStats(c *gin.Context) {
	var stats struct {
		OrderTrend []struct {
			Date       string  `json:"date"`
			OrderCount int64   `json:"order_count"`
			Amount     float64 `json:"amount"`
		} `json:"order_trend"`
		PackageStats []struct {
			PackageName string  `json:"package_name"`
			OrderCount  int64   `json:"order_count"`
			Amount      float64 `json:"amount"`
		} `json:"package_stats"`
		PaymentStats []struct {
			PaymentMethod string  `json:"payment_method"`
			OrderCount    int64   `json:"order_count"`
			Amount        float64 `json:"amount"`
		} `json:"payment_stats"`
	}

	// 获取订单趋势
	database.DB.Raw(`
		WITH dates AS (
			SELECT generate_series(
				date_trunc('day', now()) - interval '29 days',
				date_trunc('day', now()),
				interval '1 day'
			)::date AS date
		)
		SELECT 
			d.date::text,
			COUNT(DISTINCT CASE WHEN o.created_at::date = d.date THEN o.id END) as order_count,
			COALESCE(SUM(CASE WHEN o.created_at::date = d.date THEN o.amount ELSE 0 END), 0) as amount
		FROM dates d
		LEFT JOIN orders o ON TRUE
		GROUP BY d.date
		ORDER BY d.date DESC
	`).Scan(&stats.OrderTrend)

	// 获取套餐统计
	database.DB.Raw(`
		SELECT 
			p.name as package_name,
			COUNT(DISTINCT o.id) as order_count,
			COALESCE(SUM(o.amount), 0) as amount
		FROM packages p
		LEFT JOIN orders o ON o.package_id = p.id
		GROUP BY p.id, p.name
		ORDER BY amount DESC
	`).Scan(&stats.PackageStats)

	// 获取支付方式统计
	database.DB.Raw(`
		SELECT 
			payment_method,
			COUNT(DISTINCT o.id) as order_count,
			COALESCE(SUM(amount), 0) as amount
		FROM payments p
		LEFT JOIN orders o ON o.id = p.order_id
		WHERE p.status = 1
		GROUP BY payment_method
		ORDER BY amount DESC
	`).Scan(&stats.PaymentStats)

	c.JSON(http.StatusOK, model.Response{
		Code:    model.Success,
		Message: "获取成功",
		Data:    stats,
	})
} 