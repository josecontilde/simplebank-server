-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts 
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts 
WHERE id = $1 AND deleted_at IS NULL LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id DESC;

-- name: UpdateAccount :one
UPDATE accounts
SET
    owner = $2,
    balance = $3,
    currency = $4
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET
    balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :one
DELETE FROM accounts
WHERE id = $1
RETURNING *;
