-- name: CreateRefreshToken :exec
-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  $3,
  NULL
)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at;


-- name: GetRefreshToken :one
SELECT token, created_at, updated_at, user_id, expires_at, revoked_at
FROM refresh_tokens
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND revoked_at IS NULL
AND expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $2  -- același timestamp pentru ambele
WHERE token = $1;
