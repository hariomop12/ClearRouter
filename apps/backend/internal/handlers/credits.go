package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/haropm/clearrouter/apps/backend/internal/models"
	"github.com/razorpay/razorpay-go"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditsHandler struct {
	DB *gorm.DB
}

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	Amount float64 `json:"amount" binding:"required,min=1"`
}

// CreateOrder handles the creation of a new Razorpay order
func (h *CreditsHandler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Initialize Razorpay client
	client := razorpay.NewClient(os.Getenv("RAZORPAY_key_id"), os.Getenv("RAZORPAY_key_secret"))

	// Convert amount to paise (Razorpay uses smallest currency unit)
	amountInPaise := int64(req.Amount * 100)

	// Create order data
	data := map[string]interface{}{
		"amount":          amountInPaise,
		"currency":        "INR",
		"receipt":         fmt.Sprintf("rcpt_%s", uuid.New().String()[:8]),
		"payment_capture": 1,
		"notes": map[string]interface{}{
			"user_id": userID.(uuid.UUID).String(),
		},
	}

	// Create Razorpay order
	order, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating order", "details": err.Error()})
		return
	}

	// Save order in database
	payment := &models.Payment{
		UserID:          userID.(uuid.UUID),
		RazorpayOrderID: order["id"].(string),
		Amount:          req.Amount,
		Status:          "pending",
	}

	if err := h.DB.Create(payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving order"})
		return
	}

	// Return order details to client
	c.JSON(http.StatusOK, gin.H{
		"order_id": order["id"],
		"amount":   req.Amount,
		"currency": "INR",
		"key_id":   os.Getenv("RAZORPAY_key_id"),
	})
}

// RazorpayWebhookRequest represents the webhook payload from Razorpay
type RazorpayWebhookRequest struct {
	Event   string `json:"event"`
	Payload struct {
		Payment struct {
			Entity struct {
				ID      string `json:"id"`
				Amount  int64  `json:"amount"`
				Status  string `json:"status"`
				OrderID string `json:"order_id"`
			} `json:"entity"`
		} `json:"payment"`
		Order struct {
			Entity struct {
				ID     string `json:"id"`
				Amount int64  `json:"amount"`
				Notes  struct {
					UserID string `json:"user_id"`
				} `json:"notes"`
			} `json:"entity"`
		} `json:"order"`
	} `json:"payload"`
}

// AddCredits handles the Razorpay webhook for adding credits
func (h *CreditsHandler) AddCredits(c *gin.Context) {
	// Get Razorpay webhook secret from environment
	webhookSecret := os.Getenv("RAZORPAY_WEBHOOK_SECRET")

	// Verify Razorpay signature
	signature := c.GetHeader("X-Razorpay-Signature")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
		return
	}

	// Verify signature
	if !verifyWebhookSignature(string(body), signature, webhookSecret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	var webhookReq RazorpayWebhookRequest
	if err := json.Unmarshal(body, &webhookReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Only process successful payments
	if webhookReq.Event != "payment.captured" {
		c.JSON(http.StatusOK, gin.H{"message": "Event ignored"})
		return
	}

	// Convert amount to credits (amount is in paise, 1 INR = 100 paise)
	amount := float64(webhookReq.Payload.Payment.Entity.Amount) / 100
	credits := amount // 1 INR = 1 credit (you can adjust this ratio)

	// Get user ID from order notes
	userID := webhookReq.Payload.Order.Entity.Notes.UserID
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in order"})
		return
	}

	// Begin transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
		return
	}

	// Update payment status
	if err := tx.Model(&models.Payment{}).
		Where("razorpay_order_id = ?", webhookReq.Payload.Order.Entity.ID).
		Updates(map[string]interface{}{
			"razorpay_payment_id": webhookReq.Payload.Payment.Entity.ID,
			"status":              "completed",
		}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating payment"})
		return
	}

	// Update credits
	if err := tx.Exec(`
		INSERT INTO credits (user_id, total_credits, used_credits)
		VALUES (?, ?, 0)
		ON CONFLICT (user_id) 
		DO UPDATE SET total_credits = credits.total_credits + ?`,
		userID, credits, credits).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating credits"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credits added successfully"})
}

// GetCredits returns the current credit balance for a user
func (h *CreditsHandler) GetCredits(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var credits models.Credits
	if err := h.DB.Where("user_id = ?", userID.(uuid.UUID)).First(&credits).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"total_credits":     0,
				"used_credits":      0,
				"available_credits": 0,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching credits"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_credits":     credits.TotalCredits,
		"used_credits":      credits.UsedCredits,
		"available_credits": credits.TotalCredits - credits.UsedCredits,
	})
}

// Helper function to verify Razorpay webhook signature
func verifyWebhookSignature(body, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
