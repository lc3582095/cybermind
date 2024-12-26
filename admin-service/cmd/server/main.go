package main

import (
	"fmt"
	"log"

	"cybermind/admin-service/configs"
	"cybermind/admin-service/internal/api/router"
	"cybermind/admin-service/pkg/database"
)

func main() {
	// 加载配置
	config, err := configs.LoadConfig("configs/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	if err := database.InitDB(&config.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// 初始化Redis连接
	if err := database.InitRedis(&config.Redis); err != nil {
		log.Fatalf("Failed to init redis: %v", err)
	}

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.Server.Port)
	log.Printf("Server is running at %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 