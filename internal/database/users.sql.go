// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (hashed_password,email) 
VALUES (
  $1,
  $2
  )

RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type CreateUserParams struct {
	HashedPassword string
	Email          string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.HashedPassword, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const downgradeUser = `-- name: DowngradeUser :exec
UPDATE users set is_chirpy_red = false WHERE id = $1
`

func (q *Queries) DowngradeUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, downgradeUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red
FROM users
WHERE
users.email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const resetUsers = `-- name: ResetUsers :exec
DELETE FROM users
`

func (q *Queries) ResetUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetUsers)
	return err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users SET email = $1, hashed_password = $2, updated_at = NOW() WHERE id = $3

RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type UpdateUserParams struct {
	Email          string
	HashedPassword string
	ID             uuid.UUID
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser, arg.Email, arg.HashedPassword, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const upgradeUser = `-- name: UpgradeUser :exec
UPDATE users set is_chirpy_red = true WHERE id = $1
`

func (q *Queries) UpgradeUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, upgradeUser, id)
	return err
}
