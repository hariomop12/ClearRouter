package models

import "time"

// ChatMessage represents a message in a chat completion request/response
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionsRequest represents a chat completion API request
type ChatCompletionsRequest struct {
	// Optional chat to append messages to. If empty, a new chat may be created by the handler.
	ChatID   string        `json:"chat_id,omitempty"`
	Model    string        `json:"model" binding:"required"`
	Messages []ChatMessage `json:"messages" binding:"required,min=1"`
	MaxTokens *int         `json:"max_tokens,omitempty"`
	Stream    bool         `json:"stream"`
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
	ModelID      *string   `json:"model_id,omitempty" gorm:"type:uuid"` // Optional reference to models table
	Model        string    `json:"model"`                                // Model name as string
	Provider     string    `json:"provider"`                             // Provider name
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
	Cost         float64   `json:"cost" gorm:"type:numeric(12,4)"`
	Currency     string    `json:"currency" gorm:"type:text;default:'INR'"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}