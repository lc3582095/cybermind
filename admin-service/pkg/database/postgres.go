package database

import (
	"fmt"
	"log"

	"cybermind/admin-service/configs"
	"cybermind/admin-service/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(config *configs.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移数据库
	err = autoMigrate()
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	// 初始化超级管理员
	err = initSuperAdmin()
	if err != nil {
		return fmt.Errorf("failed to init super admin: %w", err)
	}

	log.Println("Database connected and initialized successfully")
	return nil
}

// autoMigrate 自动迁移数据库
func autoMigrate() error {
	return DB.AutoMigrate(
		&model.Admin{},
		&model.AdminLoginHistory{},
		&model.AdminOperation{},
		&model.Payment{},
		&model.PaymentCallback{},
		&model.PaymentRefund{},
	)
}

// initSuperAdmin 初始化超级管理员
func initSuperAdmin() error {
	var count int64
	DB.Model(&model.Admin{}).Count(&count)
	if count > 0 {
		return nil
	}

	admin := &model.Admin{
		Username:     "admin",
		Email:        "2986309418@qq.com",
		PasswordHash: "$2a$10$YF0jCx4bB4qgv2ym0wQeI.VqHw8NR0RG0Q7Q4Q8Q8Q8Q8Q8Q8Q", // 密码：admin123
		Role:         2, // 超级管理员
		Status:       1,
	}

	return DB.Create(admin).Error
} 