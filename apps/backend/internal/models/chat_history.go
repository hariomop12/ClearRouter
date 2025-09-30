package models

import (
	"time"
	"github.com/google/uuid"
)

// Chat represents a chat conversation
type Chat struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Title     string    `json:"title" gorm:"size:255;not null;default:'New Chat'"`
	Model     string    `json:"model" gorm:"size:100;not null"`
	Provider  string    `json:"provider" gorm:"size:100;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Associations
	Messages []ChatHistoryMessage `json:"messages,omitempty" gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
	User     User                 `json:"user,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

// ChatHistoryMessage represents a single message in a chat conversation stored in database
type ChatHistoryMessage struct {
	ID         uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ChatID     uuid.UUID `json:"chat_id" gorm:"type:uuid;not null"`
	Role       string    `json:"role" gorm:"size:20;not null;check:role IN ('user', 'assistant', 'system')"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	TokenCount int       `json:"token_count" gorm:"default:0"`
	Cost       float64   `json:"cost" gorm:"type:numeric(12,6);default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Associations
	Chat Chat `json:"chat,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}
// to use the existing migration table name 'chat_messages'.
func (ChatHistoryMessage) TableName() string { return "chat_messages" }

// NewChatRequest represents the request to create a new chat
type NewChatRequest struct {
	Title    string `json:"title,omitempty"`
	Model    string `json:"model" binding:"required"`
	Provider string `json:"provider,omitempty"`
}

// NewChatResponse represents the response when creating a new chat
type NewChatResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Model     string    `json:"model"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatHistoryResponse represents the response for chat history
type ChatHistoryResponse struct {
	Chats      []ChatWithMessageCount `json:"chats"`
	TotalCount int64                  `json:"total_count"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
}

// ChatWithMessageCount includes message count for chat history
type ChatWithMessageCount struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Model        string    `json:"model"`
	Provider     string    `json:"provider"`
	MessageCount int64     `json:"message_count"`
	LastMessage  string    `json:"last_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ChatDetailResponse represents a single chat with all its messages
type ChatDetailResponse struct {
	ID        uuid.UUID            `json:"id"`
	Title     string               `json:"title"`
	Model     string               `json:"model"`
	Provider  string               `json:"provider"`
	Messages  []ChatHistoryMessage `json:"messages"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}
