-- name: CreateAccount :one
INSERT INTO bank (first, last, email) VALUES ($1, $2, $3)
RETURNING *;
;

-- name: GetAccount :one
SELECT * FROM bank
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM bank
ORDER BY creation;
