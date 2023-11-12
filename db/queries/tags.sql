-- name: CreateTag :one
INSERT INTO tags (user_id, name, sign, kind)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FindTagByID :one
SELECT *
FROM tags
WHERE id = $1
  AND user_id = $2;

-- name: FindTagsByIDs :many
SELECT *
FROM tags
WHERE id = ANY (sqlc.arg(tag_ids)::BIGINT[])
  AND user_id = $1;

-- name: UpdateTagByID :one
UPDATE tags
SET name = $2,
    sign = $3,
    kind = $4
WHERE id = $1
  AND user_id = $5
RETURNING *;

-- name: DeleteAllTag :exec
DELETE
FROM tags;
