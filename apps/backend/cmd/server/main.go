package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/hariomop12/clearrouter/apps/backend/internal/handlers"
	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
	"github.com/hariomop12/clearrouter/apps/backend/internal/providers"
	"github.com/hariomop12/clearrouter/apps/backend/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	db.AutoMigrate(
		&models.User{},
		&models.APIKey{},
		&models.Credits{},
		&models.Payment{},
		&models.APIUsageLog{},
	)

	// Initialize router
	r := gin.Default()

	// Initialize providers and services
	providerService := services.NewProviderService()
	providerService.RegisterProvider(providers.NewOpenAIProvider())
	providerService.RegisterProvider(providers.NewGoogleProvider())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	apiKeyHandler := &handlers.Handler{DB: db}
	creditsHandler := &handlers.CreditsHandler{DB: db}
	chatHandler := handlers.NewChatHandler(db, providerService)

	// Routes
	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/signup", authHandler.Signup)
		auth.GET("/verify", authHandler.Verify)
		auth.POST("/login", authHandler.Login)
	}

	// Protected API Key routes
	keys := r.Group("/keys", authHandler.AuthMiddleware())
	{
		keys.POST("/create", apiKeyHandler.CreateAPIKey)
		keys.GET("", apiKeyHandler.ListAPIKeys)
	}

	// Credits routes
	credits := r.Group("/credits")
	{
		credits.POST("/order", authHandler.AuthMiddleware(), creditsHandler.CreateOrder) // Create order
		credits.POST("/add", creditsHandler.AddCredits)                                  // Razorpay webhook
		credits.GET("", authHandler.AuthMiddleware(), creditsHandler.GetCredits)         // Protected route
	}
	// Chat routes
	v1 := r.Group("/v1")
	{
		v1.POST("/chat/completions", chatHandler.ChatCompletions)
	}

	// Public routes
	r.GET("/models", handlers.GetModelsHandler)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "ClearRouter API"})
	})

	fmt.Println("Server starting on :8080")
	r.Run(":8080")
}
