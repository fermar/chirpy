-- +goose Up
CREATE TABLE users(
  id UUID primary key default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  email text not null unique

);

-- +goose Down
DROP TABLE users;
