// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package queries

import (
	"context"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (username, password_hash, salt) 
VALUES ($1, $2, $3)
RETURNING id
`

func (q *Queries) CreateAccount(ctx context.Context, username string, passwordHash string, salt string) (int32, error) {
	row := q.db.QueryRow(ctx, createAccount, username, passwordHash, salt)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getAccountID = `-- name: GetAccountID :one
SELECT id FROM accounts
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetAccountID(ctx context.Context, username string) (int32, error) {
	row := q.db.QueryRow(ctx, getAccountID, username)
	var id int32
	err := row.Scan(&id)
	return id, err
}
