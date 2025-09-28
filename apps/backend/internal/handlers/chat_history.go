package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
	"gorm.io/gorm"
)

type ChatHistoryHandler struct {
	db *gorm.DB
}

func NewChatHistoryHandler(db *gorm.DB) *ChatHistoryHandler {
	return &ChatHistoryHandler{db: db}
}

// CreateNewChat handles POST /newchat - creates a new chat conversation
func (h *ChatHistoryHandler) CreateNewChat(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	parsedUserID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse request body
	var req models.NewChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Set default title if not provided
	if strings.TrimSpace(req.Title) == "" {
		req.Title = "New Chat"
	}

	// Determine provider from model if not specified
	if req.Provider == "" {
		if strings.HasPrefix(req.Model, "gpt-") {
			req.Provider = "openai"
		} else if strings.HasPrefix(req.Model, "gemini-") {
			req.Provider = "google"
		} else {
			req.Provider = "unknown"
		}
	}

	// Create new chat
	chat := models.Chat{
		UserID:   parsedUserID,
		Title:    req.Title,
		Model:    req.Model,
		Provider: req.Provider,
	}

	if err := h.db.Create(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat", "details": err.Error()})
		return
	}

	// Return response
	response := models.NewChatResponse{
		ID:        chat.ID,
		UserID:    chat.UserID,
		Title:     chat.Title,
		Model:     chat.Model,
		Provider:  chat.Provider,
		CreatedAt: chat.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetChatHistory handles GET /chathistory - retrieves user's chat history
func (h *ChatHistoryHandler) GetChatHistory(c *gin.Context) {
	// Get user from context
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	parsedUserID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse pagination parameters
	page := 1
	pageSize := 20

	if pageParam := c.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := c.Query("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int64
	if err := h.db.Model(&models.Chat{}).Where("user_id = ?", parsedUserID).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count chats"})
		return
	}

	// Get chats with message counts and last message
	var chats []models.ChatWithMessageCount
	query := `
		SELECT 
			c.id,
			c.title,
			c.model,
			c.provider,
			c.created_at,
			c.updated_at,
			COALESCE(message_stats.message_count, 0) as message_count,
			COALESCE(message_stats.last_message, '') as last_message
		FROM chats c
		LEFT JOIN (
			SELECT 
				chat_id,
				COUNT(*) as message_count,
				COALESCE(
					(SELECT content 
					 FROM chat_messages cm2 
					 WHERE cm2.chat_id = cm.chat_id 
					 ORDER BY cm2.created_at DESC 
					 LIMIT 1), 
					''
				) as last_message
			FROM chat_messages cm
			GROUP BY chat_id
		) message_stats ON c.id = message_stats.chat_id
		WHERE c.user_id = ?
		ORDER BY c.updated_at DESC
		LIMIT ? OFFSET ?
	`

	if err := h.db.Raw(query, parsedUserID, pageSize, offset).Scan(&chats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat history", "details": err.Error()})
		return
	}

	// Return response
	response := models.ChatHistoryResponse{
		Chats:      chats,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	c.JSON(http.StatusOK, response)
}

// GetChatDetail handles GET /chathistory/:chatId - retrieves a specific chat with all messages
func (h *ChatHistoryHandler) GetChatDetail(c *gin.Context) {
	// Get user from context
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	parsedUserID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get chat ID from URL
	chatIDParam := c.Param("chatId")
	chatID, err := uuid.Parse(chatIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID format"})
		return
	}

	// Get chat with messages
	var chat models.Chat
	if err := h.db.Where("id = ? AND user_id = ?", chatID, parsedUserID).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		First(&chat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat", "details": err.Error()})
		return
	}

	// Return response
	response := models.ChatDetailResponse{
		ID:        chat.ID,
		Title:     chat.Title,
		Model:     chat.Model,
		Provider:  chat.Provider,
		Messages:  chat.Messages,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// SaveChatMessage saves a new message to an existing chat
func (h *ChatHistoryHandler) SaveChatMessage(chatID uuid.UUID, role, content string, tokenCount int, cost float64) error {
	message := models.ChatHistoryMessage{
		ChatID:     chatID,
		Role:       role,
		Content:    content,
		TokenCount: tokenCount,
		Cost:       cost,
	}

	if err := h.db.Create(&message).Error; err != nil {
		return err
	}

	// Update chat's updated_at timestamp
	return h.db.Model(&models.Chat{}).Where("id = ?", chatID).Update("updated_at", time.Now()).Error
}

// DeleteChat handles DELETE /chathistory/:chatId - deletes a chat and all its messages
func (h *ChatHistoryHandler) DeleteChat(c *gin.Context) {
	// Get user from context
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	parsedUserID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get chat ID from URL
	chatIDParam := c.Param("chatId")
	chatID, err := uuid.Parse(chatIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID format"})
		return
	}

	// Delete chat (messages will be deleted due to CASCADE)
	result := h.db.Where("id = ? AND user_id = ?", chatID, parsedUserID).Delete(&models.Chat{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat", "details": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}
