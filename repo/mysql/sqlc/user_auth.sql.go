// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: user_auth.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const createUserAuthOnDuplicateUpdateTokenAndExpiredAt = `-- name: CreateUserAuthOnDuplicateUpdateTokenAndExpiredAt :execresult
INSERT INTO user_auth (
    user_id,
    token,
    expired_at 
) VALUES (
    ?,
    ?,
    ?
)
ON DUPLICATE KEY UPDATE
    token = ?,
    expired_at = ?
`

type CreateUserAuthOnDuplicateUpdateTokenAndExpiredAtParams struct {
	UserID      []byte    `json:"user_id"`
	Token       string    `json:"token"`
	ExpiredAt   time.Time `json:"expired_at"`
	Token_2     string    `json:"token_2"`
	ExpiredAt_2 time.Time `json:"expired_at_2"`
}

func (q *Queries) CreateUserAuthOnDuplicateUpdateTokenAndExpiredAt(ctx context.Context, arg CreateUserAuthOnDuplicateUpdateTokenAndExpiredAtParams) (sql.Result, error) {
	return q.exec(ctx, q.createUserAuthOnDuplicateUpdateTokenAndExpiredAtStmt, createUserAuthOnDuplicateUpdateTokenAndExpiredAt,
		arg.UserID,
		arg.Token,
		arg.ExpiredAt,
		arg.Token_2,
		arg.ExpiredAt_2,
	)
}

const getUserAuthByUserId = `-- name: GetUserAuthByUserId :one
SELECT
    id, created_at, updated_at, user_id, token, expired_at
FROM
    user_auth
WHERE
    1 = 1
    AND user_id = ?
`

func (q *Queries) GetUserAuthByUserId(ctx context.Context, userID []byte) (UserAuth, error) {
	row := q.queryRow(ctx, q.getUserAuthByUserIdStmt, getUserAuthByUserId, userID)
	var i UserAuth
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Token,
		&i.ExpiredAt,
	)
	return i, err
}