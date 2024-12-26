package model

import (
	"time"
)

// Payment 支付记录
type Payment struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	OrderID       int64     `json:"order_id" gorm:"index"`
	PaymentNo     string    `json:"payment_no" gorm:"uniqueIndex"`
	PaymentMethod string    `json:"payment_method" gorm:"index"`
	Amount        float64   `json:"amount"`
	Status        int       `json:"status" gorm:"index"` // 0:待支付 1:支付成功 2:支付失败 3:已退款
	PaymentTime   time.Time `json:"payment_time"`
	RefundTime    time.Time `json:"refund_time"`
	CreatedAt     time.Time `json:"created_at" gorm:"index"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PaymentCallback 支付回调记录
type PaymentCallback struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	PaymentID  int64     `json:"payment_id" gorm:"index"`
	CallbackNo string    `json:"callback_no" gorm:"uniqueIndex"`
	Status     int       `json:"status" gorm:"index"` // 0:处理中 1:成功 2:失败
	CreatedAt  time.Time `json:"created_at" gorm:"index"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PaymentRefund 退款记录
type PaymentRefund struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	PaymentID int64     `json:"payment_id" gorm:"index"`
	RefundNo  string    `json:"refund_no" gorm:"uniqueIndex"`
	Amount    float64   `json:"amount"`
	Reason    string    `json:"reason"`
	Status    int       `json:"status" gorm:"index"` // 0:处理中 1:成功 2:失败
	CreatedAt time.Time `json:"created_at" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at"`
} 