package models

import (
	"time"

	"github.com/google/uuid"
)

// Credits represents a user's credit balance
type Credits struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;unique"`
	TotalCredits float64   `json:"total_credits" gorm:"type:numeric(12,2);default:0"`
	UsedCredits  float64   `json:"used_credits" gorm:"type:numeric(12,2);default:0"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
