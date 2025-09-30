package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
)

type AnalyticsHandler struct {
	DB *gorm.DB
}

func NewAnalyticsHandler(db *gorm.DB) *AnalyticsHandler {
	return &AnalyticsHandler{DB: db}
}

// GetUsageStatsWithAPIKey returns usage statistics using API key authentication
func (h *AnalyticsHandler) GetUsageStatsWithAPIKey(c *gin.Context) {
	// Get API key from header
	apiKey := c.GetHeader("Authorization")
	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
		return
	}

	// Remove "Bearer " prefix if present
	if len(apiKey) > 7 && apiKey[:7] == "Bearer " {
		apiKey = apiKey[7:]
	}

	// Find the API key and get user ID
	var key models.APIKey
	if err := h.DB.Where("api_key = ? AND active = ?", apiKey, true).First(&key).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	userID := key.UserID

	// Get query parameters
	days := c.DefaultQuery("days", "30")
	daysInt, err := strconv.Atoi(days)
	if err != nil || daysInt < 1 || daysInt > 365 {
		daysInt = 30
	}

	startDate := time.Now().AddDate(0, 0, -daysInt)

	// Get total stats
	var totalStats struct {
		TotalRequests int64   `json:"total_requests"`
		TotalTokens   int64   `json:"total_tokens"`
		TotalCost     float64 `json:"total_cost"`
	}

	h.DB.Model(&models.APIUsageAnalytics{}).
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Select("COUNT(*) as total_requests, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Scan(&totalStats)

	// Get top models
	var topModels []models.ModelUsage
	h.DB.Model(&models.APIUsageAnalytics{}).
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Select("model_requested as model, COUNT(*) as requests, SUM(total_tokens) as tokens, SUM(total_cost) as cost").
		Group("model_requested").
		Order("cost DESC").
		Limit(10).
		Scan(&topModels)

	// Get top providers
	var topProviders []models.ProviderUsage
	h.DB.Model(&models.APIUsageAnalytics{}).
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Select("provider, COUNT(*) as requests, SUM(total_tokens) as tokens, SUM(total_cost) as cost").
		Group("provider").
		Order("cost DESC").
		Limit(10).
		Scan(&topProviders)

	// Get daily breakdown
	var dailyStats []models.DailyStats
	h.DB.Model(&models.APIUsageAnalytics{}).
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Select("DATE(created_at) as date, COUNT(*) as requests, SUM(total_tokens) as tokens, SUM(total_cost) as cost").
		Group("DATE(created_at)").
		Order("date DESC").
		Scan(&dailyStats)

	stats := models.UsageStats{
		TotalRequests:  totalStats.TotalRequests,
		TotalTokens:    totalStats.TotalTokens,
		TotalCost:      totalStats.TotalCost,
		TopModels:      topModels,
		TopProviders:   topProviders,
		DailyBreakdown: dailyStats,
	}

	c.JSON(http.StatusOK, stats)
}

// GetUsageStats returns comprehensive usage statistics
func (h *AnalyticsHandler) GetUsageStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get query parameters
	days := c.DefaultQuery("days", "30")
	daysInt, err := strconv.Atoi(days)
	if err != nil || daysInt < 1 || daysInt > 365 {
		daysInt = 30
	}

	startDate := time.Now().AddDate(0, 0, -daysInt)
	userUUID := userID.(uuid.UUID)

	// Get total stats
	var totalStats struct {
		TotalRequests int64   `json:"total_requests"`
		TotalTokens   int64   `json:"total_tokens"`
		TotalCost     float64 `json:"total_cost"`
	}

	h.DB.Model(&models.APIUsageAnalytics{}).
		Where("user_id = ? AND created_at >= ?", userUUID, startDate).
		Select("COUNT(*) as total_requests, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Scan(&totalStats)

	// Get top models
	var topModels []models.ModelUsage
	h.DB.Model(&models.APIUsageAnalytics{}).
		Where("user_id = ? AND created_at >= ?", userUUID, startDate).
		Select("model_requested as model, COUNT(*) as requests, SUM(total_tokens) as tokens, SUM(total_cost) as cost").
		Group("model_requested").
		Order("cost DESC").
		Limit(10).
		Scan(&topModels)

	// Get top providers
	var topProviders []models.ProviderUsage
	h.DB.Model(&models.APIUsageAnalytics{}).
		Where("user_id = ? AND created_at >= ?", userUUID, startDate).
		Select("provider, COUNT(*) as requests, SUM(total_tokens) as tokens, SUM(total_cost) as cost").
		Group("provider").
		Order("cost DESC").
		Limit(10).
		Scan(&topProviders)

	// Get daily breakdown
	var dailyStats []models.DailyStats
	h.DB.Model(&models.APIUsageAnalytics{}).
		Where("user_id = ? AND created_at >= ?", userUUID, startDate).
		Select("DATE(created_at) as date, COUNT(*) as requests, SUM(total_tokens) as tokens, SUM(total_cost) as cost").
		Group("DATE(created_at)").
		Order("date DESC").
		Scan(&dailyStats)

	stats := models.UsageStats{
		TotalRequests:  totalStats.TotalRequests,
		TotalTokens:    totalStats.TotalTokens,
		TotalCost:      totalStats.TotalCost,
		TopModels:      topModels,
		TopProviders:   topProviders,
		DailyBreakdown: dailyStats,
	}
	c.JSON(http.StatusOK, stats)
}

// GetDailySummary returns daily usage summary
func (h *AnalyticsHandler) GetDailySummary(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    days := c.DefaultQuery("days", "30")
    daysInt, err := strconv.Atoi(days)
    if err != nil || daysInt < 1 || daysInt > 365 {
        daysInt = 30
    }

    startDate := time.Now().AddDate(0, 0, -daysInt)
    userUUID := userID.(uuid.UUID)

    // Compute daily summaries on the fly from APIUsageAnalytics
    type row struct {
        Date              time.Time `json:"date"`
        Requests          int64     `json:"requests"`
        TotalInputTokens  int64     `json:"total_input_tokens"`
        TotalOutputTokens int64     `json:"total_output_tokens"`
        TotalTokens       int64     `json:"total_tokens"`
        TotalCost         float64   `json:"total_cost"`
    }

    var results []row
    if err := h.DB.Model(&models.APIUsageAnalytics{}).
        Where("user_id = ? AND created_at >= ?", userUUID, startDate).
        Select("DATE(created_at) as date, COUNT(*) as requests, SUM(input_tokens) as total_input_tokens, SUM(output_tokens) as total_output_tokens, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
        Group("DATE(created_at)").
        Order("date DESC").
        Scan(&results).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching daily summary"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"daily_summaries": results})
}

// GetDetailedUsage returns detailed usage logs with pagination
func (h *AnalyticsHandler) GetDetailedUsage(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Pagination parameters
    page := c.DefaultQuery("page", "1")
    limit := c.DefaultQuery("limit", "50")

    pageInt, err := strconv.Atoi(page)
    if err != nil || pageInt < 1 {
        pageInt = 1
    }

    limitInt, err := strconv.Atoi(limit)
    if err != nil || limitInt < 1 || limitInt > 1000 {
        limitInt = 50
    }

    offset := (pageInt - 1) * limitInt
    userUUID := userID.(uuid.UUID)

    // Optional filters
    model := c.Query("model")
    provider := c.Query("provider")
    status := c.Query("status")

    query := h.DB.Where("user_id = ?", userUUID)

    if model != "" {
        query = query.Where("model_requested = ?", model)
    }
    if provider != "" {
        query = query.Where("provider = ?", provider)
    }
    if status != "" {
        query = query.Where("status = ?", status)
    }

    var total int64
    query.Model(&models.APIUsageAnalytics{}).Count(&total)

    var usage []models.APIUsageAnalytics
    if err := query.Order("created_at DESC").
        Offset(offset).
        Limit(limitInt).
        Find(&usage).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching usage data"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "usage": usage,
        "pagination": gin.H{
            "page":       pageInt,
            "limit":      limitInt,
            "total":      total,
            "total_pages": (total + int64(limitInt) - 1) / int64(limitInt),
        },
    })
}

// GetCostBreakdown returns cost breakdown by model and provider
func (h *AnalyticsHandler) GetCostBreakdown(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	days := c.DefaultQuery("days", "30")
	daysInt, err := strconv.Atoi(days)
	if err != nil || daysInt < 1 || daysInt > 365 {
		daysInt = 30
	}

	startDate := time.Now().AddDate(0, 0, -daysInt)
	userUUID := userID.(uuid.UUID)

	// Cost by model
	var modelCosts []struct {
		Model      string  `json:"model"`
		InputCost  float64 `json:"input_cost"`
		OutputCost float64 `json:"output_cost"`
		TotalCost  float64 `json:"total_cost"`
		Requests   int64   `json:"requests"`
	}

	h.DB.Model(&models.APIUsageAnalytics{}).
		Where("user_id = ? AND created_at >= ?", userUUID, startDate).
		Select("model_requested as model, SUM(input_cost) as input_cost, SUM(output_cost) as output_cost, SUM(total_cost) as total_cost, COUNT(*) as requests").
		Group("model_requested").
		Order("total_cost DESC").
		Scan(&modelCosts)

	// Cost by provider
	var providerCosts []struct {
		Provider   string  `json:"provider"`
		InputCost  float64 `json:"input_cost"`
		OutputCost float64 `json:"output_cost"`
		TotalCost  float64 `json:"total_cost"`
		Requests   int64   `json:"requests"`
	}

	h.DB.Model(&models.APIUsageAnalytics{}).
		Where("user_id = ? AND created_at >= ?", userUUID, startDate).
		Select("provider, SUM(input_cost) as input_cost, SUM(output_cost) as output_cost, SUM(total_cost) as total_cost, COUNT(*) as requests").
		Group("provider").
		Order("total_cost DESC").
		Scan(&providerCosts)

	c.JSON(http.StatusOK, gin.H{
		"model_costs":    modelCosts,
		"provider_costs": providerCosts,
	})
}
