-- name: CreateChirp :one
INSERT INTO chrips (body, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: ResetChrips :exec
DELETE FROM chrips;

-- name: GetAllChrips :many
SELECT * FROM chrips ORDER BY created_at ASC;

-- name: GetAllChirpsAsc :many
SELECT * FROM chrips ORDER BY created_at ASC;

-- name: GetAllChirpsDesc :many
SELECT * FROM chrips ORDER BY created_at DESC;

-- name: GetChirpByID :one
SELECT * FROM chrips WHERE id = $1;

-- name: GetChirpsByIDAsc :many
SELECT * FROM chrips ORDER BY id ASC;

-- name: GetChirpsByIDDesc :many
SELECT * FROM chrips ORDER BY id DESC;

-- name: GetChripsByUserID :many
SELECT * FROM chrips WHERE user_id = $1 ORDER BY created_at ASC;

-- name: GetChripsByUserIDDesc :many
SELECT * FROM chrips WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetChripsByUserIDAsc :many
SELECT * FROM chrips WHERE user_id = $1 ORDER BY created_at ASC;

-- name: DeleteChirp :exec
DELETE FROM chrips WHERE id = $1;