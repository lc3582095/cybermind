package main

import (
    "cybermind/auth-service/configs"
    "cybermind/auth-service/internal/api/handler"
    "cybermind/auth-service/internal/api/middleware"
    "cybermind/auth-service/internal/service"
    "cybermind/auth-service/pkg/database"
    "log"

    "github.com/gin-gonic/gin"
)

func main() {
    // 初始化数据库连接
    dbConfig := configs.GetDatabaseConfig()
    db, err := database.NewPostgresDB(dbConfig)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // 初始化用户服务
    userService := service.NewUserService(db)
    userHandler := handler.NewUserHandler(userService)

    // 设置路由
    router := gin.Default()
    v1 := router.Group("/api/v1")
    {
        auth := v1.Group("/auth")
        {
            auth.POST("/register", userHandler.Register)
            auth.POST("/login", userHandler.Login)
        }

        users := v1.Group("/users")
        users.Use(middleware.AuthMiddleware())
        {
            users.GET("/info", userHandler.GetUserInfo)
        }
    }

    // 启动服务器
    if err := router.Run(":8081"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
} 