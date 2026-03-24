-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users 
SET 
  full_name = COALESCE(sqlc.arg(full_name), full_name), 
  email = COALESCE(sqlc.arg(email), email)
WHERE username = sqlc.arg(username)
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users 
SET hashed_password = qlc.arg(hashed_password)
WHERE username = sqlc.arg(username)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;