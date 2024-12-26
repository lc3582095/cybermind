package database

import (
    "fmt"
    "log"
    "time"
    "cybermind/auth-service/configs"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewPostgresDB(config *configs.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
        config.Host,
        config.User,
        config.Password,
        config.DBName,
        config.Port,
        config.SSLMode,
    )

    log.Printf("Connecting to database: %s:%d/%s", config.Host, config.Port, config.DBName)

    // 配置GORM日志
    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    }

    // 尝试连接数据库
    db, err := gorm.Open(postgres.Open(dsn), gormConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // 获取底层的sql.DB对象
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get underlying *sql.DB: %w", err)
    }

    // 设置连接池参数
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    // 测试连接
    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    log.Printf("Successfully connected to database")
    return db, nil
} 