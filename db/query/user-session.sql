-- name: CreateUserSession :one
INSERT INTO user_sessions (
    session_id,
    user_id,
    token,
    refresh_token,
    user_agent,
    ip,
    channel,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetUserSessionByUserID :one
SELECT * FROM user_sessions
WHERE user_id = $1 LIMIT 1;

-- name: UpdateUserSession :one
UPDATE user_sessions
SET
    token = $2,
    refresh_token = $3,
    user_agent = $4,
    ip = $5,
    channel = $6,
    expires_at = $7,
    session_id = $8
WHERE
    user_id = $1
RETURNING *;
