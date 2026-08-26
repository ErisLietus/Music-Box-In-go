-- name: GetPlaylistByUser :one
SELECT * FROM playlists
WHERE user_id = $1 AND name = $2;

-- name: CreatePlaylist :one
INSERT INTO playlists (
    id,
    user_id,
    name,
    created_at,
    is_public,
    allow_collab_edits
)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    $3,
    $4
)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: DeletePlaylist :exec
DELETE FROM playlists
where id = $1;

-- name: ChangeAllowCollab :one
SELECT * FROM playlists
WHERE allow_collab_edits = $1;