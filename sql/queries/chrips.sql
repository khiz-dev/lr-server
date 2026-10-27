-- name: CreateChirp :one
INSERT INTO chrips (body, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: ResetChrips :exec
DELETE FROM chrips;