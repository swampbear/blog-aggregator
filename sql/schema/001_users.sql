-- +goose Up
CREATE TABLE users(
    id uuid DEFAULT gen_random_uuid(),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT UNIQUE,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE users;
