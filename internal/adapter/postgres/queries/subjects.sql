-- name: CreateSubject :one
INSERT INTO subjects (name, description)
VALUES ($1, $2)
RETURNING id;
