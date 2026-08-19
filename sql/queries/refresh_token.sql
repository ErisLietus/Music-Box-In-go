-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
     $1,
     now(),
     now(),
     $2,
     $3,
     NULL
)
RETURNING *;

-- name: RevokeToken :exec 
    UPDATE refresh_tokens 
    SET revoked_at = NOW(), updated_at = NOW()
    WHERE refresh_tokens.token = $1;