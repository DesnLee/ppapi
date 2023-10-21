-- name: CreateItem :one
INSERT INTO items (user_id, amount, kind, happened_at)
VALUES ($1, $2, $3, $4)
RETURNING *;
