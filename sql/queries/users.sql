-- name: GetUser :one
SELECT id, email, name, created_at, updated_at
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, email, name, created_at, updated_at
FROM users
ORDER BY created_at DESC, id DESC
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (email, name)
VALUES ($1, $2)
RETURNING id, email, name, created_at, updated_at;

-- name: CountUsers :one
SELECT COUNT(*)::bigint AS count FROM users;

-- name: UpdateUser :one
UPDATE users
SET
    email = COALESCE(sqlc.narg('email'), email),
    name = COALESCE(sqlc.narg('name'), name)
WHERE id = sqlc.arg('id')
RETURNING id, email, name, created_at, updated_at;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = $1;
