-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;


-- name: GetChirp :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: DeleteChirp :execresult
DELETE FROM chirps
WHERE id = $1
RETURNING id;

-- name: GetAllChirpsAsc :many
SELECT id, created_at, updated_at, body, user_id 
FROM chirps 
ORDER BY created_at ASC;

-- name: GetAllChirpsDesc :many
SELECT id, created_at, updated_at, body, user_id 
FROM chirps 
ORDER BY created_at DESC;

-- name: GetChirpsByAuthorIDAsc :many
SELECT id, created_at, updated_at, body, user_id 
FROM chirps 
WHERE user_id = $1 
ORDER BY created_at ASC;

-- name: GetChirpsByAuthorIDDesc :many
SELECT id, created_at, updated_at, body, user_id 
FROM chirps 
WHERE user_id = $1 
ORDER BY created_at DESC;
