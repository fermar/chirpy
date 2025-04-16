-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token,user_id,expires_at,revoked_at) 
VALUES (
  $1,
  $2,
  NOW() + INTERVAL '60 days',
  NULL
  )

RETURNING *;

-- name: ResetRefreshTokens :exec
DELETE FROM refresh_tokens;

-- name: GetRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE
token = $1 
AND expires_at > NOW() 
AND revoked_at IS NULL
LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE
token = $1; 
