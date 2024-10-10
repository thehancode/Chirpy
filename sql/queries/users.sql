-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users; 

-- name: GetAllUsers :many
SELECT *
FROM users
ORDER BY created_at ASC;

-- name: Getuser :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET 
    updated_at = NOW(),
    email = COALESCE($2, email),
    hashed_password = COALESCE($3, hashed_password)
WHERE id = $1
RETURNING *;
