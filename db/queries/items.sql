-- name: CreateItem :one
INSERT INTO items (user_id, amount, kind, happened_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FindItemByID :one
SELECT * FROM items WHERE id = $1 LIMIT 1;

-- name: DeleteAllItem :exec
DELETE FROM items;
