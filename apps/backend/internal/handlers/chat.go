package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
	"github.com/hariomop12/clearrouter/apps/backend/internal/services"
	"gorm.io/gorm"
)

type ChatHandler struct {
	DB              *gorm.DB
	ProviderService *services.ProviderService
}

func NewChatHandler(db *gorm.DB, providerService *services.ProviderService) *ChatHandler {
	return &ChatHandler{
		DB:              db,
		ProviderService: providerService,
	}
}

func (h *ChatHandler) ChatCompletions(c *gin.Context) {
	// Get API key from header
	apiKey := c.GetHeader("Authorization")
	if apiKey == "" || len(apiKey) < 8 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}
	// Remove "Bearer " prefix if present
	apiKey = apiKey[7:]

	// Find API key in database
	var key models.APIKey
	if err := h.DB.Where("api_key = ? AND active = true", apiKey).First(&key).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	// Parse request
	var req models.ChatCompletionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Determine provider from model name
	providerID := models.GetProviderFromModel(req.Model)
	if providerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid model: %s not found", req.Model)})
		return
	}

	// Get model pricing info
	modelInfo, err := h.ProviderService.GetModelInfo(providerID, req.Model)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid model: %v", err)})
		return
	}

	// Get the provider
	provider, err := h.ProviderService.GetProvider(providerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid provider: %v", err)})
		return
	}

	// Calculate input tokens
	inputTokens, err := provider.CalculateTokens(req.Messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calculating tokens"})
		return
	}

	// Calculate estimated cost
	estimatedInputCost := float64(inputTokens) * modelInfo.InputPrice

	// Check if user has enough credits
	var credits models.Credits
	if err := h.DB.Where("user_id = ?", key.UserID).First(&credits).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No credits available"})
		return
	}

	availableCredits := credits.TotalCredits - credits.UsedCredits
	if availableCredits < estimatedInputCost {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient credits"})
		return
	}

	// Create chat completion
	resp, err := provider.CreateChatCompletion(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error from provider: %v", err)})
		return
	}

	// Calculate actual costs
	inputCost, outputCost := h.ProviderService.CalculateCost(modelInfo, resp.Usage.PromptTokens, resp.Usage.CompletionTokens)
	totalCost := inputCost + outputCost

	// Begin transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
		return
	}

	// Log detailed API usage analytics
	startTime := time.Now()
	responseTime := int(time.Since(startTime).Milliseconds())
	
	usageAnalytics := models.APIUsageAnalytics{
		ID:           uuid.NewString(),
		UserID:       key.UserID.String(),
		APIKeyID:     key.ID.String(),
		RequestID:    resp.ID,
		ModelRequested: req.Model,
		ModelUsed:    req.Model, // For now, same as requested. Can be different if we do model fallbacks
		Provider:     providerID,
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
		TotalTokens:  resp.Usage.TotalTokens,
		InputCost:    inputCost,
		OutputCost:   outputCost,
		TotalCost:    totalCost,
		InputPricePerToken:  modelInfo.InputPrice,
		OutputPricePerToken: modelInfo.OutputPrice,
		Status:       "success",
		ResponseTimeMs: &responseTime,
	}

	// Also keep the old usage log for backward compatibility
	usageLog := models.APIUsageLog{
		ID:           uuid.NewString(),
		UserID:       key.UserID.String(),
		APIKeyID:     key.ID.String(),
		ModelID:      nil,
		Model:        req.Model,
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
		Cost:         totalCost,
	}

    // Create both usage logs
    if err := tx.Create(&usageLog).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging usage"})
        return
    }

    if err := tx.Create(&usageAnalytics).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging analytics"})
        return
    }

	// Update credits
	if err := tx.Model(&credits).Update("used_credits", credits.UsedCredits+totalCost).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating credits"})
		return
	}

	// Persist chat and messages
	// Determine or create chat associated with this completion
    var chatID uuid.UUID
    if req.ChatID != "" {
        // Validate provided chat belongs to the API key's user
        parsedID, err := uuid.Parse(req.ChatID)
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat_id"})
            return
        }
        var chat models.Chat
        if err := tx.Where("id = ? AND user_id = ?", parsedID, key.UserID).First(&chat).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found or not owned by user"})
            return
        }
        chatID = chat.ID
    } else {
		// Create a new chat using first user message as title (trimmed)
		title := "New Chat"
		if len(req.Messages) > 0 {
			// Prefer last user message content if roles are present
			for i := len(req.Messages) - 1; i >= 0; i-- {
				if req.Messages[i].Content != "" { // minimal safeguard
					title = req.Messages[i].Content
					break
				}
			}
			if len(title) > 60 {
				title = title[:60]
			}
		}
		newChat := models.Chat{
			UserID:   key.UserID,
			Title:    title,
			Model:        req.Model,
			Provider: providerID,
		}
		if err := tx.Create(&newChat).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating chat"})
			return
		}
		chatID = newChat.ID
	}

	// Save last user message
	var lastUserContent string
	if len(req.Messages) > 0 {
		// pick last message content
		lastUserContent = req.Messages[len(req.Messages)-1].Content
	}
	if lastUserContent != "" {
		userMsg := models.ChatHistoryMessage{
			ChatID:     chatID,
			Role:       "user",
			Content:    lastUserContent,
			TokenCount: resp.Usage.PromptTokens,
			Cost:       inputCost,
		}
		if err := tx.Create(&userMsg).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user message"})
			return
		}
	}

	// Save assistant message (from provider response)
	var assistantContent string
	if len(resp.Choices) > 0 && resp.Choices[0].Message != nil {
		assistantContent = resp.Choices[0].Message.Content
	}
	if assistantContent != "" {
		asstMsg := models.ChatHistoryMessage{
			ChatID:     chatID,
			Role:       "assistant",
			Content:    assistantContent,
			TokenCount: resp.Usage.CompletionTokens,
			Cost:       outputCost,
		}
		if err := tx.Create(&asstMsg).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving assistant message"})
			return
		}
	}

	// Update chat updated_at
	if err := tx.Model(&models.Chat{}).Where("id = ?", chatID).Update("updated_at", time.Now()).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating chat timestamp"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, resp)
}