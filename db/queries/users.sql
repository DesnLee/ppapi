-- name: CreateUser :one
INSERT INTO users (email)
VALUES ($1)
RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1;
