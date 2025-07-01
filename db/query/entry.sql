-- name: CreateEntry :one
INSERT INTO entries (
    accounts_id,
    amount
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntrys :many
SELECT * FROM entries
ORDER BY id DESC;

-- name: UpdateEntry :one
UPDATE entries
SET
    accounts_id = $2,
    amount = $3
WHERE id = $1
RETURNING *;

-- name: DeleteEntry :one
DELETE FROM entries
WHERE id = $1
RETURNING *;
