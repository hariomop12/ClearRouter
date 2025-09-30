-- migrate:up

-- Create API Usage Analytics table for detailed tracking
CREATE TABLE IF NOT EXISTS api_usage_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    api_key_id UUID NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
    
    -- Request Details
    request_id VARCHAR(255) NOT NULL,
    model_requested VARCHAR(255) NOT NULL,
    model_used VARCHAR(255) NOT NULL, -- Actual model used (after mapping)
    provider VARCHAR(100) NOT NULL,
    
    -- Token Usage
    input_tokens INTEGER NOT NULL DEFAULT 0,
    output_tokens INTEGER NOT NULL DEFAULT 0,
    total_tokens INTEGER NOT NULL DEFAULT 0,
    
    -- Cost Breakdown
    input_cost DECIMAL(15,8) NOT NULL DEFAULT 0,
    output_cost DECIMAL(15,8) NOT NULL DEFAULT 0,
    total_cost DECIMAL(15,8) NOT NULL DEFAULT 0,
    
    -- Pricing Info (for historical tracking)
    input_price_per_token DECIMAL(15,8) NOT NULL,
    output_price_per_token DECIMAL(15,8) NOT NULL,
    
    -- Request Metadata
    status VARCHAR(50) NOT NULL DEFAULT 'success', -- success, error, timeout
    error_message TEXT,
    response_time_ms INTEGER,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create Daily Usage Summary table for quick analytics
CREATE TABLE IF NOT EXISTS daily_usage_summary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    
    -- Aggregated Data
    total_requests INTEGER NOT NULL DEFAULT 0,
    total_input_tokens BIGINT NOT NULL DEFAULT 0,
    total_output_tokens BIGINT NOT NULL DEFAULT 0,
    total_tokens BIGINT NOT NULL DEFAULT 0,
    total_cost DECIMAL(15,8) NOT NULL DEFAULT 0,
    
    -- Model Breakdown (JSON for flexibility)
    models_used JSONB DEFAULT '{}',
    providers_used JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Unique constraint to prevent duplicates
    UNIQUE(user_id, date)
);

-- Create Model Pricing History table to track price changes
CREATE TABLE IF NOT EXISTS model_pricing_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_id VARCHAR(255) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    
    -- Pricing
    input_price DECIMAL(15,8) NOT NULL,
    output_price DECIMAL(15,8) NOT NULL,
    
    -- Metadata
    effective_from TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    effective_until TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_api_usage_analytics_user_id ON api_usage_analytics(user_id);
CREATE INDEX IF NOT EXISTS idx_api_usage_analytics_created_at ON api_usage_analytics(created_at);
CREATE INDEX IF NOT EXISTS idx_api_usage_analytics_model ON api_usage_analytics(model_requested);
CREATE INDEX IF NOT EXISTS idx_api_usage_analytics_provider ON api_usage_analytics(provider);
CREATE INDEX IF NOT EXISTS idx_api_usage_analytics_status ON api_usage_analytics(status);

CREATE INDEX IF NOT EXISTS idx_daily_usage_summary_user_date ON daily_usage_summary(user_id, date);
CREATE INDEX IF NOT EXISTS idx_daily_usage_summary_date ON daily_usage_summary(date);

CREATE INDEX IF NOT EXISTS idx_model_pricing_history_model ON model_pricing_history(model_id, provider);
CREATE INDEX IF NOT EXISTS idx_model_pricing_history_active ON model_pricing_history(is_active);

-- Create function to update daily summary
CREATE OR REPLACE FUNCTION update_daily_usage_summary()
RETURNS TRIGGER AS $$
BEGIN
    -- Insert or update daily summary
    INSERT INTO daily_usage_summary (
        user_id, 
        date, 
        total_requests, 
        total_input_tokens, 
        total_output_tokens, 
        total_tokens, 
        total_cost,
        models_used,
        providers_used
    )
    VALUES (
        NEW.user_id,
        DATE(NEW.created_at),
        1,
        NEW.input_tokens,
        NEW.output_tokens,
        NEW.total_tokens,
        NEW.total_cost,
        jsonb_build_object(NEW.model_requested, 1),
        jsonb_build_object(NEW.provider, 1)
    )
    ON CONFLICT (user_id, date)
    DO UPDATE SET
        total_requests = daily_usage_summary.total_requests + 1,
        total_input_tokens = daily_usage_summary.total_input_tokens + NEW.input_tokens,
        total_output_tokens = daily_usage_summary.total_output_tokens + NEW.output_tokens,
        total_tokens = daily_usage_summary.total_tokens + NEW.total_tokens,
        total_cost = daily_usage_summary.total_cost + NEW.total_cost,
        models_used = daily_usage_summary.models_used || jsonb_build_object(NEW.model_requested, 
            COALESCE((daily_usage_summary.models_used->NEW.model_requested)::integer, 0) + 1),
        providers_used = daily_usage_summary.providers_used || jsonb_build_object(NEW.provider,
            COALESCE((daily_usage_summary.providers_used->NEW.provider)::integer, 0) + 1),
        updated_at = NOW();
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update daily summary
CREATE TRIGGER trigger_update_daily_usage_summary
    AFTER INSERT ON api_usage_analytics
    FOR EACH ROW
    EXECUTE FUNCTION update_daily_usage_summary();

-- migrate:down

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_daily_usage_summary ON api_usage_analytics;
DROP FUNCTION IF EXISTS update_daily_usage_summary();

-- Drop indexes
DROP INDEX IF EXISTS idx_model_pricing_history_active;
DROP INDEX IF EXISTS idx_model_pricing_history_model;
DROP INDEX IF EXISTS idx_daily_usage_summary_date;
DROP INDEX IF EXISTS idx_daily_usage_summary_user_date;
DROP INDEX IF EXISTS idx_api_usage_analytics_status;
DROP INDEX IF EXISTS idx_api_usage_analytics_provider;
DROP INDEX IF EXISTS idx_api_usage_analytics_model;
DROP INDEX IF EXISTS idx_api_usage_analytics_created_at;
DROP INDEX IF EXISTS idx_api_usage_analytics_user_id;

-- Drop tables
DROP TABLE IF EXISTS model_pricing_history;
DROP TABLE IF EXISTS daily_usage_summary;
DROP TABLE IF EXISTS api_usage_analytics;
