package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/hariomop12/clearrouter/apps/backend/internal/dbmigrate"
	"github.com/hariomop12/clearrouter/apps/backend/internal/handlers"
	"github.com/hariomop12/clearrouter/apps/backend/internal/middleware"
	"github.com/hariomop12/clearrouter/apps/backend/internal/providers"
	"github.com/hariomop12/clearrouter/apps/backend/internal/seed"
	"github.com/hariomop12/clearrouter/apps/backend/internal/services"
	"golang.org/x/time/rate"
)

func loadEnv() {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("Unable to get working directory, skipping .env lookup")
		return
	}

	candidates := []string{
		"apps/backend/.env",
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
	}

	for _, rel := range candidates {
		path := rel
		if !filepath.IsAbs(rel) {
			path = filepath.Clean(filepath.Join(wd, rel))
		}
		if _, statErr := os.Stat(path); statErr != nil {
			continue
		}
		if loadErr := godotenv.Load(path); loadErr == nil {
			log.Printf("Loaded env file: %s", path)
			return
		}
	}

	log.Println("No .env file found, using environment variables")
}

func main() {
	loadEnv()

	fmt.Println("[STARTUP] Connecting to database...")
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("DATABASE_URL"),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to database is :", err)
	}
	fmt.Println("[STARTUP] Database connected")

	fmt.Println("[STARTUP] Running schema checks...")
	_ = dbmigrate.EnsureUsageTracking(db)

	fmt.Println("[STARTUP] Running migrations...")
	migrations := []string{
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS google_id VARCHAR(255) UNIQUE`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS github_id VARCHAR(255) UNIQUE`,
	}
	for _, m := range migrations {
		if err := db.Exec(m).Error; err != nil {
			log.Printf("[WARN] Migration failed: %v", err)
		}
	}
	fmt.Println("[STARTUP] Seeding default user...")
	seed.SeedDefaultUser(db)

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := gin.Default()

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal("Failed to set trusted proxies:", err)
	}

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:3002",
			"https://clearrouter.vercel.app",
			// Sevalla production domain(s)
			"https://clear-router-2t6fu.sevalla.app",
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow Vercel preview deployments like https://<project>-<hash>.vercel.app
			if strings.HasPrefix(origin, "https://") && strings.HasSuffix(origin, ".vercel.app") {
				return true
			}
			// Allow Sevalla app domains like https://<app>.sevalla.app
			if strings.HasPrefix(origin, "https://") && strings.HasSuffix(origin, ".sevalla.app") {
				return true
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize providers and services
	providerService := services.NewProviderService()
	providerService.RegisterProvider(providers.NewOpenAIProvider())
	providerService.RegisterProvider(providers.NewGoogleProvider())
	providerService.RegisterProvider(providers.NewAnthropicProvider())
	providerService.RegisterProvider(providers.NewDeepSeekProvider())
	providerService.RegisterProvider(providers.NewMistralProvider())

	// Rate limiter: 1 request per 20 min per user
	chatRateLimiter := middleware.NewPerUserRateLimiter(rate.Every(20*time.Minute), 1)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	oauthHandler := handlers.NewOAuthHandler(db)
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
		// convenience alias for clients using 'signin'
		auth.POST("/signin", authHandler.Login)
		// OAuth routes
		auth.GET("/google", oauthHandler.GoogleLogin)
		auth.GET("/google/callback", oauthHandler.GoogleCallback)
		auth.GET("/github", oauthHandler.GitHubLogin)
		auth.GET("/github/callback", oauthHandler.GitHubCallback)
		auth.GET("/status", oauthHandler.OAuthStatus)
	}

	// User management routes (protected)
	user := r.Group("/user", authHandler.AuthMiddleware())
	{
		user.PUT("/username", authHandler.UpdateUsername)
		user.DELETE("/account", authHandler.DeleteAccount)
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

	// Dashboard Chat routes (protected with JWT + rate limited)
	dashboardChat := r.Group("/", authHandler.AuthMiddleware())
	{
		dashboardChat.POST("/chat", chatRateLimiter.Middleware(), chatHandler.DashboardChatCompletions)
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
	r.GET("/health", healthHandler.SuperHealth)
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
