package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cybermind/admin-service/internal/api/handler"
	"cybermind/admin-service/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupBenchmarkDB() {
	var err error
	database.DB, err = gorm.Open(postgres.Open("host=dbconn.sealosbja.site port=37550 user=postgres password=wkzhx7jn dbname=postgres sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func BenchmarkGetPaymentList(b *testing.B) {
	setupBenchmarkDB()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/payments", handler.GetPaymentList)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/payments", nil)
		r.ServeHTTP(w, req)
	}
}

func BenchmarkGetPaymentDetail(b *testing.B) {
	setupBenchmarkDB()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/payments/:id", handler.GetPaymentDetail)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/payments/1", nil)
		r.ServeHTTP(w, req)
	}
}

func BenchmarkCreateRefund(b *testing.B) {
	setupBenchmarkDB()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/payments/:id/refund", handler.CreateRefund)

	reqBody := handler.CreateRefundRequest{
		Amount: 50,
		Reason: "压力测试退款",
	}
	reqData, _ := json.Marshal(reqBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/payments/1/refund", bytes.NewBuffer(reqData))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
	}
}

func BenchmarkParallelGetPaymentList(b *testing.B) {
	setupBenchmarkDB()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/payments", handler.GetPaymentList)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/payments", nil)
			r.ServeHTTP(w, req)
		}
	})
}

func BenchmarkParallelGetPaymentDetail(b *testing.B) {
	setupBenchmarkDB()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/payments/:id", handler.GetPaymentDetail)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/payments/1", nil)
			r.ServeHTTP(w, req)
		}
	})
}

func BenchmarkParallelCreateRefund(b *testing.B) {
	setupBenchmarkDB()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/payments/:id/refund", handler.CreateRefund)

	reqBody := handler.CreateRefundRequest{
		Amount: 50,
		Reason: "压力测试退款",
	}
	reqData, _ := json.Marshal(reqBody)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/payments/1/refund", bytes.NewBuffer(reqData))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
		}
	})
} 