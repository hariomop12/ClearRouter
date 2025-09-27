package models

import (
	"time"

	"github.com/google/uuid"
)

// Payment represents a payment record in the system
type Payment struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	RazorpayOrderID   string    `json:"razorpay_order_id" gorm:"type:varchar(255);not null"`
	RazorpayPaymentID string    `json:"razorpay_payment_id" gorm:"type:varchar(255)"`
	Amount            float64   `json:"amount" gorm:"type:numeric(12,2);not null"`
	Status            string    `json:"status" gorm:"type:varchar(50);default:'pending'"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
