-- name: CreateUserEmailVerificationCode :execresult
INSERT INTO user_email_verification_codes (
    user_id,
    email,
    verification_code,
    max_try,
    expired_at 
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
)
;

-- name: GetUserEmailVerificationCodeByEmailAndVerificationCode :one
SELECT
    *
FROM
    user_email_verification_codes
WHERE
    1 = 1
    AND email = ?
    AND verification_code = ?
;

-- name: DecreaseUserEmailVerificationCodeMaxTryById :exec
UPDATE
    user_email_verification_codes
SET
    max_try = max_try - 1
WHERE
    id = ? 
;

-- name: DeleteUserEmailVerificationCodeByUserId :exec
DELETE FROM
    user_email_verification_codes
WHERE
    user_id = ? 
;

-- name: GetLastCreatedUserEmailVerificationCodeByUserId :one
SELECT
    *
FROM
    user_email_verification_codes
WHERE
    1 = 1
    AND user_id = ? 
ORDER BY
    created_at DESC
LIMIT
    1
;