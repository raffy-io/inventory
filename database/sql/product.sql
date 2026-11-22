-- name: ListProducts :many
SELECT * FROM products;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (name,price,quantity)
VALUES
    ($1,$2,$3)
RETURNING *;

-- name: UpdateProduct :exec
UPDATE products
SET
    name = $2,
    price = $3,
    quantity =$4
WHERE 
    id = $1; 

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
