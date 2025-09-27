package models

import (
	"time"

	"github.com/google/uuid"
)

// APIKey represents an API key in the system
// APIKey represents an API key in the system
type APIKey struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	APIKey    string    `json:"api_key" gorm:"column:api_key;type:text;not null;unique"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:now()"`
}

// TableName specifies the table name for GORM
func (APIKey) TableName() string {
	return "api_keys"
}
