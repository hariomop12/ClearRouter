package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONB type for PostgreSQL JSONB fields
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

// APIUsageAnalytics represents detailed API usage tracking
type APIUsageAnalytics struct {
	ID       string  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID   string  `json:"user_id" gorm:"type:uuid;not null"`
	APIKeyID *string `json:"api_key_id" gorm:"type:uuid"`

	RequestID      string `json:"request_id" gorm:"not null"`
	ModelRequested string `json:"model_requested" gorm:"not null"`
	ModelUsed      string `json:"model_used" gorm:"not null"`
	Provider       string `json:"provider" gorm:"not null"`

	// Token Usage
	InputTokens  int `json:"input_tokens" gorm:"default:0"`
	OutputTokens int `json:"output_tokens" gorm:"default:0"`
	TotalTokens  int `json:"total_tokens" gorm:"default:0"`

	// Cost Breakdown
	InputCost  float64 `json:"input_cost" gorm:"type:decimal(15,8);default:0"`
	OutputCost float64 `json:"output_cost" gorm:"type:decimal(15,8);default:0"`
	TotalCost  float64 `json:"total_cost" gorm:"type:decimal(15,8);default:0"`
	Currency   string  `json:"currency" gorm:"type:text;default:'INR'"`

	// Pricing Info (for historical tracking)
	InputPricePerToken  float64 `json:"input_price_per_token" gorm:"type:decimal(15,8)"`
	OutputPricePerToken float64 `json:"output_price_per_token" gorm:"type:decimal(15,8)"`

	Status         string  `json:"status" gorm:"default:success"`
	ErrorMessage   *string `json:"error_message,omitempty"`
	ResponseTimeMs *int    `json:"response_time_ms,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// DailyUsageSummary represents aggregated daily usage data
type DailyUsageSummary struct {
	ID   string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Date time.Time `json:"date" gorm:"type:date;not null"`

	// Aggregated Data
	TotalRequests     int     `json:"total_requests" gorm:"default:0"`
	TotalInputTokens  int64   `json:"total_input_tokens" gorm:"default:0"`
	TotalOutputTokens int64   `json:"total_output_tokens" gorm:"default:0"`
	TotalTokens       int64   `json:"total_tokens" gorm:"default:0"`
	TotalCost         float64 `json:"total_cost" gorm:"type:decimal(15,8);default:0"`

	// Model and Provider Breakdown
	ModelsUsed    JSONB `json:"models_used" gorm:"type:jsonb;default:'{}'"`
	ProvidersUsed JSONB `json:"providers_used" gorm:"type:jsonb;default:'{}'"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// ModelPricingHistory tracks pricing changes over time
type ModelPricingHistory struct {
	ID       string `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ModelID  string `json:"model_id" gorm:"not null"`
	Provider string `json:"provider" gorm:"not null"`

	// Pricing
	InputPrice  float64 `json:"input_price" gorm:"type:decimal(15,8);not null"`
	OutputPrice float64 `json:"output_price" gorm:"type:decimal(15,8);not null"`

	// Metadata
	EffectiveFrom  time.Time  `json:"effective_from" gorm:"default:now()"`
	EffectiveUntil *time.Time `json:"effective_until,omitempty"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// UsageStats represents usage statistics for analytics
type UsageStats struct {
	TotalRequests  int64           `json:"total_requests"`
	TotalTokens    int64           `json:"total_tokens"`
	TotalCost      float64         `json:"total_cost"`
	TopModels      []ModelUsage    `json:"top_models"`
	TopProviders   []ProviderUsage `json:"top_providers"`
	DailyBreakdown []DailyStats    `json:"daily_breakdown"`
}

type ModelUsage struct {
	Model    string  `json:"model"`
	Requests int64   `json:"requests"`
	Tokens   int64   `json:"tokens"`
	Cost     float64 `json:"cost"`
}

type ProviderUsage struct {
	Provider string  `json:"provider"`
	Requests int64   `json:"requests"`
	Tokens   int64   `json:"tokens"`
	Cost     float64 `json:"cost"`
}

type DailyStats struct {
	Date     string  `json:"date"`
	Requests int64   `json:"requests"`
	Tokens   int64   `json:"tokens"`
	Cost     float64 `json:"cost"`
}
