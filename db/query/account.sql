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

-- PASSING THREE ARGUMENTS: owner, balance, currency
-- RETURNING * WILL RETURN THE VALUE OF ALL COLUMNS AS WE WANT ID OF THE ACCOUNT ONCE CREATED

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;   -- FOR UPDATE LOCKS THE ROW FOR UPDATE ; DO NOT UPDATE THE KEY - WILL STOP DEADLOCKS

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;
