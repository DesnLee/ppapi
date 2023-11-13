-- name: CreateTag :one
INSERT INTO tags (user_id, name, sign, kind)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FindTagByID :one
SELECT *
FROM tags
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL;

-- name: FindTagsByIDs :many
SELECT *
FROM tags
WHERE id = ANY (sqlc.arg(tag_ids)::BIGINT[])
  AND user_id = $1
  AND deleted_at IS NULL;

-- name: UpdateTagByID :one
UPDATE tags
SET name = COALESCE(NULLIF(@name, ''), name),
    sign = COALESCE(NULLIF(@sign, ''), sign),
    kind = COALESCE(NULLIF(@kind, ''), kind)
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteTagByID :exec
UPDATE tags
SET deleted_at = NOW()
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL;

-- name: ListTagsByUserIDWithCondition :many
SELECT *
FROM tags
WHERE user_id = $1
  AND kind = $2
  AND deleted_at IS NULL
ORDER BY created_at DESC
OFFSET $3 LIMIT $4;

-- name: CountTagsByUserIDWithCondition :one
SELECT COUNT(*)
FROM tags
WHERE user_id = $1
  AND kind = $2
  AND deleted_at IS NULL;

-- name: DeleteAllTag :exec
DELETE
FROM tags;
