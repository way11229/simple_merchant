// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: users.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :execresult
INSERT INTO users (
    name,
    email,
    password
) VALUES (
    ?,
    ?,
    ?
)
`

type CreateUserParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.exec(ctx, q.createUserStmt, createUser, arg.Name, arg.Email, arg.Password)
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
    id, created_at, updated_at, name, email, email_verified_at, password
FROM
    users
WHERE
    1 = 1
    AND email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.getUserByEmailStmt, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.EmailVerifiedAt,
		&i.Password,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT
    id, created_at, updated_at, name, email, email_verified_at, password
FROM
    users
WHERE
    1 = 1
    AND id = ?
`

func (q *Queries) GetUserById(ctx context.Context, id int32) (User, error) {
	row := q.queryRow(ctx, q.getUserByIdStmt, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.EmailVerifiedAt,
		&i.Password,
	)
	return i, err
}
