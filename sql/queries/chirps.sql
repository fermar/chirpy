-- name: CreateChirp :one
INSERT INTO chirps (body,user_id) 
VALUES (
  $1,
  $2
  )

RETURNING *;

-- name: DeleteChirpByID :exec
DELETE FROM chirps WHERE id = $1;

-- name: GetAllChirps :many
SELECT *
FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT *
FROM chirps
WHERE
id = $1;


-- name: GetChirpByUserID :many
SELECT *
FROM chirps
WHERE
user_id = $1
ORDER BY created_at ASC;
