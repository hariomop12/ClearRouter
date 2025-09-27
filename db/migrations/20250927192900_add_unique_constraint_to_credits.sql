-- migrate:up
ALTER TABLE credits ADD CONSTRAINT unique_user_id UNIQUE (user_id);

-- migrate:down
ALTER TABLE credits DROP CONSTRAINT unique_user_id;