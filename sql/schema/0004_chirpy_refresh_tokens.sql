-- +goose Up
CREATE TABLE refresh_tokens(
  token text primary key,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  user_id  uuid not null,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  expires_at timestamp not null,
  revoked_at timestamp 

);

-- +goose Down
DROP TABLE refresh_tokens;
