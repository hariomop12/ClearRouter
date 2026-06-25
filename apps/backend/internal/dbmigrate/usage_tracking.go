package dbmigrate

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// EnsureUsageTracking makes the minimal, idempotent schema changes needed for
// usage logging and analytics. This keeps the API working even when the DB was
// initialized from an older schema.sql.
func EnsureUsageTracking(db *gorm.DB) error {
	statements := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`,

		`ALTER TABLE api_usage_logs ADD COLUMN IF NOT EXISTS model VARCHAR(255);`,
		`ALTER TABLE api_usage_logs ADD COLUMN IF NOT EXISTS provider VARCHAR(100);`,
		`ALTER TABLE api_usage_logs ADD COLUMN IF NOT EXISTS currency VARCHAR(10) DEFAULT 'INR';`,

		`CREATE TABLE IF NOT EXISTS api_usage_analytics (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			api_key_id UUID NULL REFERENCES api_keys(id) ON DELETE CASCADE,
			request_id VARCHAR(255) NOT NULL,
			model_requested VARCHAR(255) NOT NULL,
			model_used VARCHAR(255) NOT NULL,
			provider VARCHAR(100) NOT NULL,
			input_tokens INTEGER NOT NULL DEFAULT 0,
			output_tokens INTEGER NOT NULL DEFAULT 0,
			total_tokens INTEGER NOT NULL DEFAULT 0,
			input_cost DECIMAL(15,8) NOT NULL DEFAULT 0,
			output_cost DECIMAL(15,8) NOT NULL DEFAULT 0,
			total_cost DECIMAL(15,8) NOT NULL DEFAULT 0,
			currency VARCHAR(10) DEFAULT 'INR',
			input_price_per_token DECIMAL(15,8),
			output_price_per_token DECIMAL(15,8),
			status VARCHAR(50) NOT NULL DEFAULT 'success',
			response_time_ms INTEGER,
			error_message TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,

		`CREATE INDEX IF NOT EXISTS idx_api_usage_analytics_user_id ON api_usage_analytics(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_api_usage_analytics_created_at ON api_usage_analytics(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_api_usage_analytics_model_requested ON api_usage_analytics(model_requested);`,
		`CREATE INDEX IF NOT EXISTS idx_api_usage_analytics_provider ON api_usage_analytics(provider);`,
	}

	for i, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			msg := err.Error()
			if strings.Contains(msg, "duplicate") || strings.Contains(msg, "already exists") {
				continue
			}
			return fmt.Errorf("usage tracking schema step %d failed: %w", i, err)
		}
	}

	return nil
}
