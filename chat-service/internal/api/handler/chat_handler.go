package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"cybermind/chat-service/internal/service"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		chatService: &service.ChatService{},
	}
}

// CreateConversation 创建对话
func (h *ChatHandler) CreateConversation(c *gin.Context) {
	var req struct {
		ModelID int64  `json:"model_id" binding:"required"`
		Title   string `json:"title"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1001, "message": "参数错误", "error": err.Error()})
		return
	}

	// 从上下文获取用户ID(假设已经通过中间件设置)
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	conversation, err := h.chatService.CreateConversation(userIDInt, req.ModelID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1005, "message": "创建对话失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": conversation})
}

// GetConversation 获取对话详情
func (h *ChatHandler) GetConversation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1001, "message": "参数错误", "error": err.Error()})
		return
	}

	conversation, err := h.chatService.GetConversation(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1004, "message": "对话不存在", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": conversation})
}

// ListConversations 获取对话列表
func (h *ChatHandler) ListConversations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	conversations, total, err := h.chatService.ListConversations(userIDInt, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1005, "message": "获取对话列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"total": total,
			"list":  conversations,
		},
	})
}

// AddMessage 添加消息
func (h *ChatHandler) AddMessage(c *gin.Context) {
	var req struct {
		ConversationID int64  `json:"conversation_id" binding:"required"`
		Role           string `json:"role" binding:"required"`
		Content        string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1001, "message": "参数错误", "error": err.Error()})
		return
	}

	message, err := h.chatService.AddMessage(req.ConversationID, req.Role, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1005, "message": "添加消息失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": message})
}

// GetMessages 获取消息列表
func (h *ChatHandler) GetMessages(c *gin.Context) {
	conversationID, err := strconv.ParseInt(c.Param("conversation_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1001, "message": "参数错误", "error": err.Error()})
		return
	}

	messages, err := h.chatService.GetMessages(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1005, "message": "获取消息列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": messages})
} 