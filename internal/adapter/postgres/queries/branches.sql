-- name: CreateBranch :one
INSERT INTO branches (name, address, city, status)
VALUES ($1, $2, $3, 'ACTIVE')
RETURNING id;
