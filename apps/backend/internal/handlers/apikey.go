package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
	"gorm.io/gorm"
)

// Handler struct holds dependencies for the API handlers
type Handler struct {
	DB *gorm.DB
}

// CreateAPIKey handles the creation of a new API key for a user
func (h *Handler) CreateAPIKey(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Generate a secure random API key
	key, err := generateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating API key"})
		return
	}

	// Create new API key in database
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type", "actual_type": fmt.Sprintf("%T", userID)})
		return
	}

	// Create API key model
	apiKey := &models.APIKey{
		UserID: userUUID,
		APIKey: key,
		Active: true,
	}

	// Debug information
	fmt.Printf("Generated key: %s\n", key)
	fmt.Printf("Creating API key model: %+v\n", apiKey)

	// Debug info
	fmt.Printf("Creating API key with data: %+v\n", apiKey)

	// Attempt to create with debug mode
	result := h.DB.Debug().Create(apiKey)
	if err := result.Error; err != nil {
		fmt.Printf("Error creating API key: %v\nRows affected: %d\n", err, result.RowsAffected)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "Error creating API key",
			"details":      err.Error(),
			"api_key_data": apiKey,
			"sql_error":    fmt.Sprintf("%+v", err),
		})
		return
	}

	// Print success information
	fmt.Printf("Successfully created API key. ID: %s\n", apiKey.ID)

	c.JSON(http.StatusCreated, apiKey)
}

// ListAPIKeys handles listing all API keys for a user
func (h *Handler) ListAPIKeys(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var apiKeys []models.APIKey
	if err := h.DB.Where("user_id = ?", userID.(uuid.UUID)).Order("created_at DESC").Find(&apiKeys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching API keys"})
		return
	}

	c.JSON(http.StatusOK, apiKeys)
}

// generateAPIKey creates a secure random API key
func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
