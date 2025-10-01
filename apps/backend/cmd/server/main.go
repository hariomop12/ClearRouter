package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/hariomop12/clearrouter/apps/backend/internal/handlers"
	"github.com/hariomop12/clearrouter/apps/backend/internal/providers"
	"github.com/hariomop12/clearrouter/apps/backend/internal/seed"
	"github.com/hariomop12/clearrouter/apps/backend/internal/services"
)

func main() {
	// Load .env file if it exists (development mode)
	// In production, environment variables are provided by Docker
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Database connection
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Schema is managed by SQL migrations in /db/migrations/
	// No AutoMigrate needed here

	// Seed default user (idempotent)
	seed.SeedDefaultUser(db)

	// Initialize router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize providers and services
	providerService := services.NewProviderService()
	providerService.RegisterProvider(providers.NewOpenAIProvider())
	providerService.RegisterProvider(providers.NewGoogleProvider())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	apiKeyHandler := handlers.NewHandler(db)
	creditsHandler := handlers.NewCreditsHandler(db)
	chatHandler := handlers.NewChatHandler(db, providerService)
	chatHistoryHandler := handlers.NewChatHistoryHandler(db)
	analyticsHandler := handlers.NewAnalyticsHandler(db)
	healthHandler := handlers.NewHealthHandler(db)

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
		keys.DELETE("/:id", apiKeyHandler.DeleteAPIKey)
	}

	// Credits routes
	credits := r.Group("/credits")
	{
		credits.POST("/order", authHandler.AuthMiddleware(), creditsHandler.CreateOrder)    // Create order
		credits.POST("/verify", authHandler.AuthMiddleware(), creditsHandler.VerifyPayment) // Verify payment
		credits.POST("/add", creditsHandler.AddCredits)                                     // Razorpay webhook
		credits.GET("", authHandler.AuthMiddleware(), creditsHandler.GetCredits)            // Protected route
	}
	// Chat routes
	v1 := r.Group("/v1")
	{
		v1.POST("/chat/completions", chatHandler.ChatCompletions)
	}

	// Dashboard Chat routes (protected with JWT)
	dashboardChat := r.Group("/", authHandler.AuthMiddleware())
	{
		dashboardChat.POST("/chat", chatHandler.DashboardChatCompletions)
	}

	// Chat History routes (protected)
	chatHistory := r.Group("/", authHandler.AuthMiddleware())
	{
		chatHistory.POST("/newchat", chatHistoryHandler.CreateNewChat)
		chatHistory.GET("/chathistory", chatHistoryHandler.GetChatHistory)
		chatHistory.GET("/chathistory/:chatId", chatHistoryHandler.GetChatDetail)
		chatHistory.DELETE("/chathistory/:chatId", chatHistoryHandler.DeleteChat)
	}

	// Analytics routes (protected)
	analytics := r.Group("/analytics", authHandler.AuthMiddleware())
	{
		analytics.GET("/usage", analyticsHandler.GetUsageStats)
		analytics.GET("/daily", analyticsHandler.GetDailySummary)
		analytics.GET("/detailed", analyticsHandler.GetDetailedUsage)
		analytics.GET("/costs", analyticsHandler.GetCostBreakdown)
	}

	// Public routes
	r.GET("/models", handlers.GetModelsHandler)
	r.GET("/health/super", healthHandler.SuperHealth)

	// API Key-based analytics (no JWT required)
	r.GET("/api/usage", analyticsHandler.GetUsageStatsWithAPIKey)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "ClearRouter API"})
	})

	fmt.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
