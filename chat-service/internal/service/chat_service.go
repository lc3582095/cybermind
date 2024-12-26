package service

import (
	"errors"
	"cybermind/chat-service/internal/model"
	"cybermind/chat-service/pkg/database"
)

type ChatService struct{}

// CreateConversation 创建新对话
func (s *ChatService) CreateConversation(userID, modelID int64, title string) (*model.Conversation, error) {
	conversation := &model.Conversation{
		UserID:  userID,
		ModelID: modelID,
		Title:   title,
	}

	if err := database.DB.Create(conversation).Error; err != nil {
		return nil, err
	}
	return conversation, nil
}

// GetConversation 获取对话详情
func (s *ChatService) GetConversation(id int64) (*model.Conversation, error) {
	var conversation model.Conversation
	if err := database.DB.Preload("Messages").First(&conversation, id).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}

// ListConversations 获取用户的对话列表
func (s *ChatService) ListConversations(userID int64, page, size int) ([]model.Conversation, int64, error) {
	var conversations []model.Conversation
	var total int64

	offset := (page - 1) * size
	query := database.DB.Model(&model.Conversation{}).Where("user_id = ?", userID)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&conversations).Error; err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// AddMessage 添加消息
func (s *ChatService) AddMessage(conversationID int64, role, content string) (*model.Message, error) {
	// 检查对话是否存在
	var conversation model.Conversation
	if err := database.DB.First(&conversation, conversationID).Error; err != nil {
		return nil, errors.New("对话不存在")
	}

	message := &model.Message{
		ConversationID: conversationID,
		Role:          role,
		Content:       content,
	}

	if err := database.DB.Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// GetMessages 获取对话消息列表
func (s *ChatService) GetMessages(conversationID int64) ([]model.Message, error) {
	var messages []model.Message
	if err := database.DB.Where("conversation_id = ?", conversationID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// UpdateConversationPoints 更新对话消耗的积分
func (s *ChatService) UpdateConversationPoints(conversationID int64, points int) error {
	return database.DB.Model(&model.Conversation{}).Where("id = ?", conversationID).
		UpdateColumn("points_consumed", database.DB.Raw("points_consumed + ?", points)).Error
} 