package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
	"github.com/razorpay/razorpay-go"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditsHandler struct {
	DB *gorm.DB
}

func firstEnv(keys ...string) string {
	for _, k := range keys {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return ""
}

// NewCreditsHandler creates a new CreditsHandler instance
func NewCreditsHandler(db *gorm.DB) *CreditsHandler {
	return &CreditsHandler{DB: db}
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

	keyID := firstEnv("RAZORPAY_KEY_ID", "RAZORPAY_key_id")
	keySecret := firstEnv("RAZORPAY_KEY_SECRET", "RAZORPAY_key_secret")
	if keyID == "" || keySecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay credentials are not configured"})
		return
	}

	// Initialize Razorpay client
	client := razorpay.NewClient(keyID, keySecret)

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
		// preferred fields (frontend)
		"id":       order["id"],
		"amount":   amountInPaise,
		"currency": "INR",
		"key":      keyID,

		// backwards-compatible aliases
		"order_id": order["id"],
		"key_id":   keyID,
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
	webhookSecret := firstEnv("RAZORPAY_WEBHOOK_SECRET")

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

// VerifyPayment handles payment verification from frontend
func (h *CreditsHandler) VerifyPayment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
		RazorpayOrderID   string `json:"razorpay_order_id"`
		RazorpaySignature string `json:"razorpay_signature"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Payment verification binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	log.Printf("Payment verification request - Payment ID: %s, Order ID: %s, Signature: %s",
		req.RazorpayPaymentID, req.RazorpayOrderID, req.RazorpaySignature)

	keyID := firstEnv("RAZORPAY_KEY_ID", "RAZORPAY_key_id")
	keySecret := firstEnv("RAZORPAY_KEY_SECRET", "RAZORPAY_key_secret")

	// Check if this is a test mode payment (payment ID starts with pay_test)
	isTestPayment := strings.HasPrefix(req.RazorpayPaymentID, "pay_test") || strings.HasPrefix(keyID, "rzp_test")

	if keyID == "" || keySecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay credentials are not configured"})
		return
	}

	// Always prefer order_id from Razorpay (prevents mismatches if multiple pending orders exist).
	client := razorpay.NewClient(keyID, keySecret)
	shouldFetchPayment := req.RazorpayOrderID == "" || req.RazorpaySignature != ""
	if shouldFetchPayment {
		paymentEntity, err := client.Payment.Fetch(req.RazorpayPaymentID, nil, nil)
		if err != nil {
			log.Printf("Unable to fetch payment details from Razorpay: %v", err)
			if req.RazorpayOrderID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing razorpay_order_id"})
				return
			}
		} else {
			if oid, ok := paymentEntity["order_id"].(string); ok && oid != "" {
				req.RazorpayOrderID = oid
			}
			if status, ok := paymentEntity["status"].(string); ok && status != "" && status != "captured" && status != "authorized" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Payment not successful", "status": status})
				return
			}
		}
	}

	// For test payments, we might not always get order_id and/or signature; fall back to DB only if needed.
	if isTestPayment && req.RazorpayOrderID == "" {
		log.Printf("Test mode payment detected, attempting to find order by payment ID")
		var payment models.Payment
		if err := h.DB.Where("user_id = ? AND status = 'pending'", userID.(uuid.UUID)).Order("created_at DESC").First(&payment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No pending payment found"})
			return
		}
		req.RazorpayOrderID = payment.RazorpayOrderID
		log.Printf("Found matching order: %s", req.RazorpayOrderID)
	}

	if req.RazorpayOrderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing razorpay_order_id"})
		return
	}

	// Verify signature (skip for test mode only if signature is missing)
	if !isTestPayment || req.RazorpaySignature != "" {
		generatedSignature := generateSignature(req.RazorpayOrderID, req.RazorpayPaymentID, keySecret)
		if generatedSignature != req.RazorpaySignature {
			log.Printf("Signature verification failed - Expected: %s, Got: %s", generatedSignature, req.RazorpaySignature)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment signature"})
			return
		}
		log.Printf("Signature verification successful")
	} else {
		log.Printf("Skipping signature verification for test mode payment")
	}

	// Update payment status
	var payment models.Payment
	if err := h.DB.Where("razorpay_order_id = ? AND user_id = ?", req.RazorpayOrderID, userID.(uuid.UUID)).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if payment.Status == "completed" {
		c.JSON(http.StatusOK, gin.H{"message": "Payment already verified"})
		return
	}

	// Begin transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
		return
	}

	// Update payment status
	if err := tx.Model(&payment).Updates(map[string]interface{}{
		"razorpay_payment_id": req.RazorpayPaymentID,
		"status":              "completed",
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating payment"})
		return
	}

	// Add credits
	credits := payment.Amount // 1 INR = 1 credit (adjust ratio as needed)
	if err := tx.Exec(`
		INSERT INTO credits (user_id, total_credits, used_credits)
		VALUES (?, ?, 0)
		ON CONFLICT (user_id) 
		DO UPDATE SET total_credits = credits.total_credits + ?`,
		userID.(uuid.UUID), credits, credits).Error; err != nil {
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

	c.JSON(http.StatusOK, gin.H{"message": "Payment verified and credits added successfully"})
}

// Helper function to generate Razorpay signature for verification
func generateSignature(orderID, paymentID, secret string) string {
	message := orderID + "|" + paymentID
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// Helper function to verify Razorpay webhook signature
func verifyWebhookSignature(body, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
