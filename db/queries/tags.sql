-- name: CreateTag :one
INSERT INTO tags (user_id, name, sign, kind)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FindTagByID :one
SELECT *
FROM tags
WHERE id = $1;

-- name: FindTagsByIDs :many
SELECT *
FROM tags
WHERE id = ANY($1::BIGINT[]);

-- name: DeleteAllTag :exec
DELETE FROM tags;
