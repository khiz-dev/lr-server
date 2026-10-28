-- name: CreateChirp :one
INSERT INTO chrips (body, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: ResetChrips :exec
DELETE FROM chrips;

-- name: GetAllChrips :many
SELECT * FROM chrips ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT * FROM chrips WHERE id = $1;

-- name: GetChripsByUserID :many
SELECT * FROM chrips WHERE user_id = $1 ORDER BY created_at ASC;