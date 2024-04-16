-- name: CreateUserAuth :exec
INSERT INTO user_auth (
    user_id,
    token,
    expired_at 
) VALUES (
    ?,
    ?,
    ?
)
;

-- name: GetUserAuthByUserId :one
SELECT
    *
FROM
    user_auth
WHERE
    user_id = ? 
;

-- name: UpdateUserAuthById :exec
UPDATE
    user_auth
SET
    token = ?,
    expired_at = ?,
    updated_at = NOW()
WHERE
    id = ?
;

-- name: DeleteUserAuthByUserId :exec
DELETE FROM
   user_auth 
WHERE
    user_id = ? 
;