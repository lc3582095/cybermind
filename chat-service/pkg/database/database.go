package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"cybermind/chat-service/internal/model"
)

var DB *gorm.DB

func InitDB() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		"dbconn.sealosbja.site",
		"postgres",
		"wkzhx7jn",
		"postgres",
		37550,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(
		&model.Conversation{},
		&model.Message{},
	); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	DB = db
	return nil
} 