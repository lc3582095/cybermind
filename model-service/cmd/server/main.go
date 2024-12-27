package main

import (
	"log"

	"cybermind/model-service/internal/api/router"
	"cybermind/model-service/pkg/database"
)

func main() {
	// 初始化数据库连接
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 设置路由
	r := router.SetupRouter(db)

	// 启动服务器
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
