package service

import (
	"errors"
	"log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"cybermind/model-service/internal/model"
)

type ModelService struct {
	db *gorm.DB
}

func NewModelService(db *gorm.DB) *ModelService {
	return &ModelService{db: db}
}

// CreateModel 创建新的模型配置
func (s *ModelService) CreateModel(m *model.Model) error {
	log.Printf("Creating model: %+v", m)
	err := s.db.Create(m).Error
	if err != nil {
		log.Printf("Error creating model: %v", err)
	}
	return err
}

// GetModel 获取模型配置
func (s *ModelService) GetModel(id int64) (*model.Model, error) {
	log.Printf("Getting model with ID: %d", id)
	var m model.Model
	err := s.db.First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Model not found with ID: %d", id)
			return nil, errors.New("model not found")
		}
		log.Printf("Error getting model: %v", err)
		return nil, err
	}
	return &m, nil
}

// UpdateModel 更新模型配置
func (s *ModelService) UpdateModel(m *model.Model) error {
	log.Printf("Updating model: %+v", m)
	err := s.db.Save(m).Error
	if err != nil {
		log.Printf("Error updating model: %v", err)
	}
	return err
}

// ListModels 获取模型列表
func (s *ModelService) ListModels(c *gin.Context) ([]model.Model, error) {
	log.Println("Listing all models")
	var models []model.Model
	err := s.db.Find(&models).Error
	if err != nil {
		log.Printf("Error listing models: %v", err)
		return nil, err
	}
	log.Printf("Found %d models", len(models))
	for i, m := range models {
		log.Printf("Model %d: %+v", i+1, m)
	}
	return models, nil
}

// DeleteModel 删除模型
func (s *ModelService) DeleteModel(id int64) error {
	log.Printf("Deleting model with ID: %d", id)
	err := s.db.Delete(&model.Model{}, id).Error
	if err != nil {
		log.Printf("Error deleting model: %v", err)
	}
	return err
}

// UpdateModelStatus 更新模型状态
func (s *ModelService) UpdateModelStatus(id int64, status int) error {
	log.Printf("Updating model status: ID=%d, status=%d", id, status)
	err := s.db.Model(&model.Model{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		log.Printf("Error updating model status: %v", err)
	}
	return err
} 