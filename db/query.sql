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

-- name: GetProduct :one
SELECT * FROM products WHERE id = $1 LIMIT 1;

-- name: CreateInvoice :one
INSERT INTO invoices (id, user_id, invoice_number, total) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListInvoices :many
SELECT * FROM invoices WHERE user_id = $1;

-- name: CreateInvoiceItem :one
INSERT INTO invoice_items (id, user_id, product_id, invoice_id, amount) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListInvoiceItems :many
SELECT * FROM invoice_items WHERE user_id = $1 AND invoice_id = $2;
