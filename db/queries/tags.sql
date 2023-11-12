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
SET name = COALESCE(NULLIF(@name, ''), name),
    sign = COALESCE(NULLIF(@sign, ''), sign),
    kind = COALESCE(NULLIF(@kind, ''), kind)
WHERE id = $1
  AND user_id = $2
RETURNING *;

-- name: DeleteAllTag :exec
DELETE
FROM tags;
