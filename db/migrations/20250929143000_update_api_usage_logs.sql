-- migrate:up
-- Add new fields to api_usage_logs table
ALTER TABLE api_usage_logs 
ADD COLUMN IF NOT EXISTS model VARCHAR(255),
ADD COLUMN IF NOT EXISTS provider VARCHAR(100);

-- Make model_id nullable (remove NOT NULL constraint if it exists)
ALTER TABLE api_usage_logs 
ALTER COLUMN model_id DROP NOT NULL;

-- migrate:down
-- Remove the added columns
ALTER TABLE api_usage_logs 
DROP COLUMN IF EXISTS model,
DROP COLUMN IF EXISTS provider;

-- Restore model_id NOT NULL constraint
ALTER TABLE api_usage_logs 
ALTER COLUMN model_id SET NOT NULL;
