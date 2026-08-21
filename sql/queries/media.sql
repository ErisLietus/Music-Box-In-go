-- name: GetMediaByPlaylist :many
SELECT 
    users.username,
    media.title,
    media.id,
    media.type,
    media.file_url,
    media.position,
    media.created_at
FROM media
INNER JOIN users ON media.added_by_user_id = users.id
WHERE media.playlist_id = $1
ORDER BY media.position ASC;

-- name: GetMaxMediaPosition :one
SELECT COALESCE(MAX(position), 0)::INTEGER AS max_position
FROM media
WHERE playlist_id = $1;

-- name: CreateMedia :one
INSERT INTO media (
    id,
    playlist_id,
    title,
    file_url,
    created_at,
    type,
    position,
    added_by_user_id
)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW(),
    $4,
    $5,
    $6
)
RETURNING *;