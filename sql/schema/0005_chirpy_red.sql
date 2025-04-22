-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red boolean DEFAULT false;

-- Ensure all existing rows comply with the NOT NULL constraint
UPDATE users
SET  is_chirpy_red = false
WHERE  is_chirpy_red IS NULL;

-- Add the NOT NULL constraint
ALTER TABLE users
ALTER COLUMN is_chirpy_red SET NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_chirpy_red;

