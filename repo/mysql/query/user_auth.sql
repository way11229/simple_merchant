-- name: CreateUserAuthOnDuplicateUpdateTokenAndExpiredAt :execresult
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
;

-- name: GetUserAuthByUserId :one
SELECT
    *
FROM
    user_auth
WHERE
    1 = 1
    AND user_id = ? 
;

