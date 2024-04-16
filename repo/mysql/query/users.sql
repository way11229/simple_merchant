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
    1 = 1
    AND email = ?
;

-- name: GetUserById :one
SELECT
    *
FROM
    users
WHERE
    1 = 1
    AND id = ?
;

-- name: DeleteUserById :exec
DELETE FROM
   users 
WHERE
    id = ? 
;