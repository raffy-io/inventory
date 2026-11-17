-- +goose Up
CREATE TABLE items (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    price NUMERIC NOT NULL CHECK( price > 0),
    weight NUMERIC NOT NULL  CHECK( weight > 0)
);

-- +goose Down
DROP TABLE items;
