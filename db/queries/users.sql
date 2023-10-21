-- name: CreateUser :one
INSERT INTO users (email)
VALUES ($1)
RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: FindUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: DeleteAllUser :exec
DELETE FROM users;
