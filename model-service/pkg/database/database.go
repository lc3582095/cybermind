package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"cybermind/model-service/internal/model"
)

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	dsn := "host=cybermind-postgresql.ns-han88ija.svc user=postgres password=wkzhx7jn dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	log.Printf("正在连接数据库，DSN: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	log.Println("数据库连接成功")

	return db, nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate(db *gorm.DB) error {
	log.Println("开始数据库迁移...")

	// 迁移模型表、供应商表和API Key池表
	if err := db.AutoMigrate(&model.Model{}, &model.Provider{}, &model.APIKeyPool{}); err != nil {
		return err
	}
	log.Println("数据库迁移完成")

	// 初始化默认模型数据
	var modelCount int64
	if err := db.Model(&model.Model{}).Count(&modelCount).Error; err != nil {
		return err
	}

	if modelCount == 0 {
		log.Println("未找到模型数据，创建默认模型...")
		defaultModels := model.DefaultModels()
		if err := db.Create(&defaultModels).Error; err != nil {
			return err
		}
		log.Println("默认模型创建成功")
	} else {
		log.Printf("发现 %d 个现有模型", modelCount)
	}

	// 从现有模型中提取供应商信息并初始化供应商表
	var models []model.Model
	if err := db.Model(&model.Model{}).Order("created_at asc").Find(&models).Error; err != nil {
		return err
	}

	// 使用map来按供应商分组，同时记录每个供应商的最早创建时间
	providerMap := make(map[string]struct {
		models    []model.Model
		firstTime time.Time
	})

	for _, m := range models {
		data, exists := providerMap[m.Provider]
		if !exists {
			data.firstTime = m.CreatedAt
		}
		data.models = append(data.models, m)
		providerMap[m.Provider] = data
	}

	log.Printf("发现 %d 个不同的供应商", len(providerMap))

	// 为每个供应商创建记录
	for providerCode, data := range providerMap {
		var count int64
		// 检查供应商是否已存在
		if err := db.Model(&model.Provider{}).Where("code = ?", providerCode).Count(&count).Error; err != nil {
			return err
		}

		if count == 0 {
			// 创建供应商记录
			provider := model.Provider{
				Name:       providerCode,
				Code:       providerCode,
				ModelCount: len(data.models),
				SortID:     1000,
				Status:     1,
				CreatedAt:  data.firstTime,
				UpdatedAt:  data.firstTime,
			}

			if err := db.Create(&provider).Error; err != nil {
				return err
			}
			log.Printf("创建供应商记录: %s, 模型数量: %d, 创建时间: %s", providerCode, len(data.models), data.firstTime)
		}
	}

	return nil
}
