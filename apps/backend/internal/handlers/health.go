package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler provides system health checks
 type HealthHandler struct {
	 db    *gorm.DB
	 start time.Time
 }

 func NewHealthHandler(db *gorm.DB) *HealthHandler {
	 return &HealthHandler{db: db, start: time.Now()}
 }

 // SuperHealth returns an aggregated health report
 func (h *HealthHandler) SuperHealth(c *gin.Context) {
	 // Context with timeout for external checks
	 ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	 defer cancel()

	 // Check database connectivity
	 dbOK := false
	 dbErr := ""
	 if h.db != nil {
		 if sqlDB, err := h.db.DB(); err == nil {
			 if err := sqlDB.PingContext(ctx); err == nil {
				 dbOK = true
			 } else {
				 dbErr = err.Error()
			 }
		 } else {
			 dbErr = err.Error()
		 }
	 } else {
		 dbErr = "gorm DB is nil"
	 }

	 // Environment checks (presence only)
	 env := map[string]any{
		 "DATABASE_URL":   requiredEnv("DATABASE_URL"),
		 "OPENAI_API_KEY": optionalEnv("OPENAI_API_KEY"),
		 "GOOGLE_API_KEY": optionalEnv("GOOGLE_API_KEY"),
	 }

	 // Build response
	 resp := gin.H{
		 "status":       overallStatus(dbOK),
		 "timestamp":    time.Now().UTC().Format(time.RFC3339),
		 "uptime_secs":  int(time.Since(h.start).Seconds()),
		 "version":      os.Getenv("APP_VERSION"),
		 "dependencies": gin.H{
			 "postgres": gin.H{
				 "ok":    dbOK,
				 "error": dbErr,
			 },
		 },
		 "env": env,
	 }

	 code := http.StatusOK
	 if !dbOK {
		 code = http.StatusServiceUnavailable
	 }
	 c.JSON(code, resp)
 }

 func overallStatus(dbOK bool) string {
	 if dbOK {
		 return "ok"
	 }
	 return "degraded"
 }

 func requiredEnv(key string) gin.H {
	 v := os.Getenv(key)
	 return gin.H{
		 "present": v != "",
	 }
 }

 func optionalEnv(key string) gin.H {
	 v := os.Getenv(key)
	 return gin.H{
		 "present": v != "",
	 }
 }
