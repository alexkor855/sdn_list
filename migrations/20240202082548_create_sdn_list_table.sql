-- +goose Up
-- +goose StatementBegin
CREATE TABLE sdn_list (
    id SERIAL PRIMARY KEY,
    uid BIGINT NOT NULL CHECK (uid > 0),
    first_name VARCHAR,
    last_name VARCHAR,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sdn_list;
-- +goose StatementEnd