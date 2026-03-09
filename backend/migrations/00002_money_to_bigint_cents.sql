-- +goose Up
ALTER TABLE products
    ALTER COLUMN price TYPE BIGINT USING ROUND(price * 100)::BIGINT;

ALTER TABLE orders
    ALTER COLUMN total_price TYPE BIGINT USING ROUND(total_price * 100)::BIGINT;

-- +goose Down
ALTER TABLE orders
    ALTER COLUMN total_price TYPE NUMERIC USING (total_price::NUMERIC / 100);

ALTER TABLE products
    ALTER COLUMN price TYPE NUMERIC USING (price::NUMERIC / 100);

