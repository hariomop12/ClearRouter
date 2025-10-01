-- migrate:up
ALTER TABLE api_usage_logs 
ADD COLUMN model VARCHAR(255),
ADD COLUMN provider VARCHAR(100),
ADD COLUMN currency VARCHAR(10) DEFAULT 'INR';

CREATE TABLE api_usage_analytics (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id UUID NOT NULL,
  api_key_id UUID NOT NULL,
  request_id VARCHAR(255) NOT NULL,
  model_requested VARCHAR(255) NOT NULL,
  model_used VARCHAR(255) NOT NULL,
  provider VARCHAR(100) NOT NULL,
  input_tokens INTEGER DEFAULT 0,
  output_tokens INTEGER DEFAULT 0,
  total_tokens INTEGER DEFAULT 0,
  input_cost DECIMAL(15,8) DEFAULT 0,
  output_cost DECIMAL(15,8) DEFAULT 0,
  total_cost DECIMAL(15,8) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'INR',
  input_price_per_token DECIMAL(15,8),
  output_price_per_token DECIMAL(15,8),
  status VARCHAR(50) DEFAULT 'success',
  response_time_ms INTEGER,
  error_message TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- migrate:down

