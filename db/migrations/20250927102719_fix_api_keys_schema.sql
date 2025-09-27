-- migrate:up
-- First, drop the duplicate columns and recreate correctly
ALTER TABLE api_keys DROP COLUMN IF EXISTS key;

-- migrate:down
-- No down migration needed as this is a fix

