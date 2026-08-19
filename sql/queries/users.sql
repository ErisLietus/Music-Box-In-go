-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
     now(),
     now(),
     $1,
     $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users; 

-- name: CheckUserByEmail :one 
SELECT * from users
where email = $1;

-- name: GetUserFromRefreshToken :one
SELECT * from users
JOIN refresh_tokens on users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1 and refresh_tokens.expires_at > NOW() and refresh_tokens.revoked_at IS NULL;

-- name: UpdateUser :one 
UPDATE users
SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE $1 = id
RETURNING *;

-- name: UpgradeToRed :one
UPDATE users
set is_chirpy_red = TRUE
where $1 = id
RETURNING *; 