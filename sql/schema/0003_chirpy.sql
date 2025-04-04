-- +goose Up
ALTER TABLE users
ADD COLUMN hashed_password TEXT DEFAULT 'unset';

-- Ensure all existing rows comply with the NOT NULL constraint
UPDATE users
SET hashed_password = 'unset'
WHERE hashed_password IS NULL;

-- Add the NOT NULL constraint
ALTER TABLE users
ALTER COLUMN hashed_password SET NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN hashed_password;

