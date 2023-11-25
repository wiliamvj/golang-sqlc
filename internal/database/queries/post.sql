-- name: CreateOnePost :exec
INSERT INTO posts (id, title, body, author_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);
