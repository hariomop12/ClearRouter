package handlers
 

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	// Auto-migrate test tables here
	return db
}

func TestAuthHandler_Login(t *testing.T) {
	db := setupTestDB()
	handler := NewAuthHandler(db)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/login", handler.Login)

	// Test invalid credentials
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}
	jsonData, _ := json.Marshal(loginData)
	
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAPIKeyHandler_CreateAPIKey(t *testing.T) {
	// Test API key creation
	db := setupTestDB()
	handler := NewHandler(db)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/keys/create", handler.CreateAPIKey)
	
	// Add test implementation
}

func TestChatHandler_ChatCompletions(t *testing.T) {
	// Test chat completions endpoint
	// Mock provider responses
	// Test cost calculations
	// Test usage analytics logging
}