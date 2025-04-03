-- +goose Up
CREATE TABLE chirps(
  id UUID primary key default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  body text not null,
  user_id  uuid not null,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE

);

-- +goose Down
DROP TABLE chirps;
