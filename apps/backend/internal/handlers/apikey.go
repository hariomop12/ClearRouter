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

// NewHandler creates a new Handler instance
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
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

// DeleteAPIKey handles deleting an API key
func (h *Handler) DeleteAPIKey(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get API key ID from URL parameter
	keyID := c.Param("id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API key ID is required"})
		return
	}

	// Parse UUID
	keyUUID, err := uuid.Parse(keyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID format"})
		return
	}

	// Find the API key to ensure it belongs to the user
	var apiKey models.APIKey
	if err := h.DB.Where("id = ? AND user_id = ?", keyUUID, userID.(uuid.UUID)).First(&apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding API key"})
		return
	}

	// Delete the API key
	if err := h.DB.Delete(&apiKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting API key"})
		return
	}

	fmt.Printf("Successfully deleted API key. ID: %s\n", apiKey.ID)
	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
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
