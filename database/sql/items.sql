-- name: ListItems :many
SELECT * FROM items;

-- name: CreateItem :one
INSERT INTO items (name,price,weight)
VALUES
    ($1,$2,$3)
RETURNING *;

-- name: UpdateItem :exec
UPDATE items
SET name = $2,
    price = $3,
    weight = $4
WHERE id = $1;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = $1;