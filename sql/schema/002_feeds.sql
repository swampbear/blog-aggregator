-- +goose Up
CREATE TABLE feeds (
    id uuid DEFAULT gen_random_uuid(),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT UNIQUE,
    url TEXT,
    user_id uuid REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY(id)
);


-- +goose Down
DROP TABLE feeds;
