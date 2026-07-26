-- name: ListItems :many
SELECT * FROM items
ORDER BY created_at DESC, id DESC;

-- name: CreateItem :one
INSERT INTO items (
    title, url, note
) VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateItems :copyfrom
INSERT INTO items (
    title, url, note
) VALUES ($1, $2, $3);
