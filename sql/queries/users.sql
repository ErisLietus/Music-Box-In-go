-- name: CreateUser :one
INSERT INTO users (
    id,
    created_at,
    username,
    hashed_email,
    hashed_password
)
VALUES (
    gen_random_uuid(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: CheckUserByEmail :one
SELECT * FROM users
WHERE hashed_email = $1;

-- name: CheckUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserFromRefreshToken :one
SELECT * from users
JOIN refresh_tokens on users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1 and refresh_tokens.expires_at > NOW() and refresh_tokens.revoked_at IS NULL;

-- name: UpdateUser :one 
UPDATE users
SET hashed_email = $2, hashed_password = $3, updated_at = NOW()
WHERE $1 = id
RETURNING *;

-- name: UpgradeToRed :one
UPDATE users
set is_chirpy_red = TRUE
where $1 = id
RETURNING *; 