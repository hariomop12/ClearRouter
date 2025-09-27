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

	"github.com/haropm/clearrouter/apps/backend/internal/handlers"
	"github.com/haropm/clearrouter/apps/backend/internal/models"
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
	// This only migrates the User model. The SQL migration file is more complete.
	// It's better to run the SQL migration file manually.
	// You could also add all your model structs here.
	db.AutoMigrate(&models.User{}, &models.APIKey{}) // Migrate both User and APIKey models

	// Initialize router
	r := gin.Default()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	apiKeyHandler := &handlers.Handler{DB: db}

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

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "ClearRouter API"})
	})

	fmt.Println("Server starting on :8080")
	r.Run(":8080")
}
