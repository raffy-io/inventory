-- +goose Up
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    price NUMERIC NOT NULL,
    quantity NUMERIC NOT NULL
);

-- +goose Down
DROP TABLE products;
