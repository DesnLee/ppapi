-- name: CreateItem :one
INSERT INTO items (user_id, amount, kind, happened_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FindItemByID :one
SELECT *
FROM items
WHERE id = $1
LIMIT 1;

-- name: ListItemsByUserIDWithCondition :many
SELECT items.*,
       array_agg(items_tags.tag_id)::BIGINT[] AS tag_ids
FROM items
         LEFT JOIN
     items_tags ON items.id = items_tags.item_id
WHERE items.user_id = $1
  AND sqlc.arg(happened_after) <= happened_at
  AND happened_at < sqlc.arg(happened_before)
GROUP BY items.id
ORDER BY happened_at DESC
OFFSET $2 LIMIT $3;

-- name: CountItemsByUserIDWithCondition :one
SELECT COUNT(*)
FROM items
WHERE user_id = $1
  AND sqlc.arg(happened_after) <= happened_at
  AND happened_at < sqlc.arg(happened_before);

-- name: ListItemsByUserIDWithDateRange :many
SELECT items.*
FROM items
WHERE items.user_id = $1
  AND sqlc.arg(happened_after) <= happened_at
  AND happened_at < sqlc.arg(happened_before);

-- name: DeleteAllItem :exec
DELETE
FROM items;
