-- name: GetAccountID :one
SELECT id FROM accounts
WHERE username = $1 LIMIT 1;

-- name: CreateAccount :one
INSERT INTO accounts (username, password_hash, salt) 
VALUES ($1, $2, $3)
RETURNING id;

-- name: AccountExists :one
SELECT EXISTS(SELECT 1 FROM accounts WHERE username = $1);

-- name: GetAccount :one
SELECT id, password_hash, salt FROM accounts
WHERE username = $1 LIMIT 1;
