package models

import "time"

// ChatMessage represents a message in a chat completion request/response
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionsRequest represents a chat completion API request
type ChatCompletionsRequest struct {
	Model     string        `json:"model" binding:"required"`
	Messages  []ChatMessage `json:"messages" binding:"required,min=1"`
	MaxTokens *int          `json:"max_tokens,omitempty"`
	Stream    bool          `json:"stream"`
}

// ChatCompletionsResponse represents a chat completion API response
type ChatCompletionsResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int          `json:"index"`
		Message *ChatMessage `json:"message,omitempty"`
		Delta   *ChatMessage `json:"delta,omitempty"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// APIUsageLog represents an API usage log entry
type APIUsageLog struct {
	ID           string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       string    `json:"user_id" gorm:"type:uuid"`
	APIKeyID     string    `json:"api_key_id" gorm:"type:uuid"`
	Model        string    `json:"model"`
	Provider     string    `json:"provider"`
	TokensUsed   int       `json:"tokens_used"`
	Cost         float64   `json:"cost"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	Status       string    `json:"status"`
	RequestID    string    `json:"request_id"`
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
	InputCost    float64   `json:"input_cost"`
	OutputCost   float64   `json:"output_cost"`
	TotalCost    float64   `json:"total_cost"`
}
