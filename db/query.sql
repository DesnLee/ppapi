-- name: CreateUser :one
INSERT INTO users (email, phone)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email   = $1,
    phone   = $2,
    address = $3
WHERE id = $4
RETURNING *;
