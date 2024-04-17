-- name: CreateUser :execresult
INSERT INTO users (
    name,
    email,
    password
) VALUES (
    ?,
    ?,
    ?
)
;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = ?
;

-- name: GetUserById :one
SELECT
    *
FROM
    users
WHERE
    id = ?
;

-- name: DeleteUserById :exec
DELETE FROM
   users 
WHERE
    id = ? 
;

-- name: VerifyUserEmailById :exec
UPDATE
    users
SET
    updated_at = NOW(),
    email_verified_at = NOW()
WHERE
    id = ? 
;