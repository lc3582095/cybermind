package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"cybermind/chat-service/internal/api/router"
	"cybermind/chat-service/pkg/database"
)

func main() {
	// 设置gin模式
	gin.SetMode(gin.DebugMode)

	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 创建gin引擎
	engine := gin.Default()

	// 注册路由
	router.RegisterRoutes(engine)
	log.Println("路由注册完成")

	// 启动服务器
	log.Println("服务器启动在 :8081")
	if err := engine.Run(":8081"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
} 