package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"cybermind/model-service/internal/api/handler"
	"cybermind/model-service/internal/api/router"
	"cybermind/model-service/internal/service"
	"cybermind/model-service/internal/model"
)

func main() {
	// 初始化数据库连接
	dsn := "host=cybermind-postgresql.ns-han88ija.svc user=postgres password=wkzhx7jn dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	log.Printf("Connecting to database with DSN: %s", dsn)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	// 自动迁移数据库表结构
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(&model.Model{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrations completed")

	// 初始化默认模型数据
	var count int64
	if err := db.Model(&model.Model{}).Count(&count).Error; err != nil {
		log.Fatalf("Failed to count models: %v", err)
	}

	if count == 0 {
		log.Println("No models found, creating default models...")
		defaultModels := model.DefaultModels()
		if err := db.Create(&defaultModels).Error; err != nil {
			log.Fatalf("Failed to create default models: %v", err)
		}
		log.Println("Default models created successfully")
	} else {
		log.Printf("Found %d existing models", count)
	}

	// 初始化服务
	modelService := service.NewModelService(db)
	modelHandler := handler.NewModelHandler(modelService)

	// 初始化Gin路由
	r := gin.Default()

	// 设置路由
	router.SetupRouter(r, modelHandler)

	// 启动服务器，监听所有网络接口
	log.Println("Server starting on :8081...")
	if err := r.Run("0.0.0.0:8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 