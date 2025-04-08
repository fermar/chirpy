-- name: CreateUser :one
INSERT INTO users (hashed_password,email) 
VALUES (
  $1,
  $2
  )

RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE
users.email = $1
LIMIT 1;
