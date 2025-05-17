-- name: CreateUser :one
INSERT INTO users (id, name, email) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;

-- name: CreateProduct :one
INSERT INTO products (id, user_id, name, price) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListProducts :many
SELECT * FROM products WHERE user_id = $1;
