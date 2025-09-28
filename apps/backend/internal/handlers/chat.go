package handlers

import (
	"fmt"
	"net/http"

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

	// Log API usage
	usageLog := models.APIUsageLog{
		ID:           uuid.NewString(),
		UserID:       key.UserID.String(),
		APIKeyID:     key.ID.String(),
		Model:        req.Model,
		Provider:     providerID,
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
		InputCost:    inputCost,
		OutputCost:   outputCost,
		TotalCost:    inputCost + outputCost,
		Status:       "success",
		RequestID:    resp.ID,
	}

	if err := tx.Create(&usageLog).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging usage"})
		return
	}

	// Update credits
	if err := tx.Model(&credits).Update("used_credits", credits.UsedCredits+totalCost).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating credits"})
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
